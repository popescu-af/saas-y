package connection

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/mock"
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
func (c *ChannelListenerMock) ProcessMessage(msg *Message, write WriteOnChannelFunc) error {
	defer c.wg.Done()

	c.Called()
	c.Msg = msg
	return nil
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
	c.self.closed = true
	c.self.cv.Signal()
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
	if c.closed {
		return nil, errChannelClosed
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

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
	if c.closed {
		return errChannelClosed
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.msgs = append(c.msgs, m)
	c.cv.Signal()

	return nil
}

// NewChannelMockEndpoint creates a channel endpoint instance.
func NewChannelMockEndpoint() *ChannelMockEndpoint {
	result := new(ChannelMockEndpoint)
	result.cv = sync.NewCond(&result.mutex)
	return result
}
