package connection

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ChannelListenerMock is a mock for the ChannelListener interface.
type ChannelListenerMock struct {
	mock.Mock
	wg  sync.WaitGroup
	Msg *Message
}

// ExpectMessageArrival implements the method with the same name from ChannelListener.
func (c *ChannelListenerMock) ExpectMessageArrival() {
	c.wg.Add(1)
	c.On("ProcessMessage")
}

// WaitAndAssertMessageArrived implements the method with the same name from ChannelListener.
func (c *ChannelListenerMock) WaitAndAssertMessageArrived(t *testing.T) {
	c.wg.Wait()
	c.AssertCalled(t, "ProcessMessage")
}

// ProcessMessage implements the method with the same name from ChannelListener.
func (c *ChannelListenerMock) ProcessMessage(msg *Message, write WriteOnChannelFunc) {
	defer c.wg.Done()

	c.Called()
	c.Msg = msg
}

// NewChannelListenerMock creates a ChannelListenerMock instance.
func NewChannelListenerMock() *ChannelListenerMock {
	return &ChannelListenerMock{}
}

// ChannelMock is a mock for the Channel interface.
type ChannelMock struct {
	self  *ChannelMockEndpoint
	other *ChannelMockEndpoint
}

// Read implements the method with the same name from Channel.
func (c *ChannelMock) Read() (*Message, error) {
	return c.self.ReadMessage()
}

// Write implements the method with the same name from Channel.
func (c *ChannelMock) Write(m *Message) error {
	return c.other.WriteMessage(m)
}

// Close implements the method with the same name from Channel.
func (c *ChannelMock) Close() {
	c.other.WriteMessage(&Message{
		Type: CloseMessage,
	})
	c.self.Close()
}

// NewChannelMock creates a ChannelMock instance.
func NewChannelMock(self, other *ChannelMockEndpoint) Channel {
	return &ChannelMock{
		self:  self,
		other: other,
	}
}

// ChannelMockEndpoint is the type for a mock channel's endpoint.
type ChannelMockEndpoint struct {
	mutex  sync.Mutex
	cv     *sync.Cond
	msgs   []*Message
	closed bool
}

var errChannelClosed = fmt.Errorf("channel is closed")

// ReadMessage reads a message from the mock channel.
func (c *ChannelMockEndpoint) ReadMessage() (*Message, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.closed {
		return nil, errChannelClosed
	}

	if len(c.msgs) == 0 {
		c.cv.Wait()
		if c.closed {
			return nil, errChannelClosed
		}
	}

	msg := c.msgs[0]
	if msg.Type == CloseMessage {
		c.closed = true
	}

	c.msgs = c.msgs[1:]
	return msg, nil
}

// WriteMessage writes a message to the mock channel.
func (c *ChannelMockEndpoint) WriteMessage(m *Message) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.closed {
		return errChannelClosed
	}

	c.msgs = append(c.msgs, m)
	c.cv.Signal()

	return nil
}

// Close closes the endpoint.
func (c *ChannelMockEndpoint) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.closed = true
	c.cv.Signal()
}

// NewChannelMockEndpoint creates a channel endpoint instance.
func NewChannelMockEndpoint() *ChannelMockEndpoint {
	result := new(ChannelMockEndpoint)
	result.cv = sync.NewCond(&result.mutex)
	return result
}

// EndpointMockInstance is useful for mocking full-duplex endpoints that
// can talk to the service to be tested.
type EndpointMockInstance struct {
	Channel  Channel
	Listener *ChannelListenerMock
	Conn     *FullDuplex
}

// ChannelListenerCreator creates channel listeners for the service side, for testing.
type ChannelListenerCreator func(t *testing.T) (ChannelListener, error)

// SpawnClientInstances spawns instances of clients of the service to be tested.
func SpawnClientInstances(t *testing.T, clientCount int, listenerCreator ChannelListenerCreator, clientName, serviceName string) ([]*EndpointMockInstance, *sync.WaitGroup) {
	var result []*EndpointMockInstance
	wg := &sync.WaitGroup{}

	for i := 0; i < clientCount; i++ {
		wg.Add(2)

		serverEndpoint := NewChannelMockEndpoint()
		clientEndpoint := NewChannelMockEndpoint()

		// Instantiate server stuff
		serverChannel := NewChannelMock(serverEndpoint, clientEndpoint)
		serverListener, err := listenerCreator(t)
		require.NoError(t, err, "unexpected failure when creating server listener")

		server := NewFullDuplex(serverListener, serverChannel, serviceName+" -> "+clientName+" : "+strconv.Itoa(i))

		go func() {
			defer wg.Done()
			defer server.Close()

			server.Run()
		}()

		// Instantiate client stuff
		clientChannel := NewChannelMock(clientEndpoint, serverEndpoint)
		clientListener := NewChannelListenerMock()
		client := &EndpointMockInstance{
			Channel:  clientChannel,
			Listener: clientListener,
			Conn:     NewFullDuplex(clientListener, clientChannel, clientName+" -> "+serviceName+" : "+strconv.Itoa(i)),
		}

		go func() {
			defer wg.Done()

			client.Conn.Run()
		}()

		result = append(result, client)
	}

	// Wait for the connections to be created
	time.Sleep(time.Duration(50) * time.Millisecond)

	return result, wg
}
