package connection

import (
	"fmt"
	"sync"
	"time"

	"github.com/popescu-af/saas-y/pkg/log"
)

// ReadFn is the type of function listening for messages on a channel.
type ReadFn func() (*Message, error)

// WriteFn is the type of function sending messages on a channel.
type WriteFn func(*Message) error

// CloseFn is the type of function that closes a channel.
type CloseFn func()

// Channel is the structure containing a channel's read/write functions.
type Channel struct {
	Read  ReadFn
	Write WriteFn
	Close CloseFn
}

// Message types, same as websocket.
const (
    // TextMessage denotes a text data message.
    TextMessage = 1

    // BinaryMessage denotes a binary data message.
    BinaryMessage = 2

    // CloseMessage denotes a close control message.
    CloseMessage = 8

    // PingMessage denotes a ping control message.
    PingMessage = 9

    // PongMessage denotes a pong control message.
    PongMessage = 10
)

// Message is a structure encompassing a message's type and payload.
type Message struct {
	Type    int
	Payload []byte
}

// FullDuplexEndpoint is the type for communication endpoints.
// Communication endpoints function in a full-duplex way, being possible to
// receive and send messages on the channel at the same time.
// Processing received messages has priority. If there is nothing to process,
// the endpoint can do housekeeping work or even send messages to the other side
// of the communication channel.
type FullDuplexEndpoint interface {
	// ProcessMessage should process the given message and react accordingly,
	// by changing state and/or writing something back or none of them.
	ProcessMessage(*Message, WriteFn) error
	// Poll should inspect the current state and decide whether to take any action
	// and/or write something on the channel or not.
	Poll(time.Time, WriteFn) error
	// CloseCommandChannel should return the channel the full duplex connection would listen
	// to for 'close' commands that come from inside the endpoint's logic.
	// A value written on the channel (either true or false) will close the connection.
	CloseCommandChannel() chan bool
}

// FullDuplex is a full-duplex connection that takes a full-duplex endpoint
// and a channel. It handles messages arriving on the channel while periodically
// polling the endpoint for new messages to be dispached on the channel.
type FullDuplex struct {
	endpoint        FullDuplexEndpoint
	channel         *Channel
	pollingPeriod   time.Duration
	stopWaitClosing chan bool
	stopReading     chan bool
	stopProcessing  chan bool
	wgStop          sync.WaitGroup
	lockState       sync.Mutex
	lockStop        sync.Mutex
	isRunning       bool
	isClosed        bool
}

// NewFullDuplex creates a new, inactive full-duplex connection.
// Call Run to run it.
func NewFullDuplex(endpoint FullDuplexEndpoint, channel *Channel, pollingPeriod time.Duration) *FullDuplex {
	return &FullDuplex{
		endpoint:        endpoint,
		channel:         channel,
		pollingPeriod:   pollingPeriod,
		stopWaitClosing: make(chan bool),
		stopReading:     make(chan bool),
		stopProcessing:  make(chan bool),
	}
}

// Run is a blocking function that handles messages arriving on the channel
// while periodically polling the endpoint for new messages to be dispached on the channel.
func (f *FullDuplex) Run() error {
	f.lockState.Lock()
	if f.isClosed {
		f.lockState.Unlock()
		return fmt.Errorf("already closed")
	}

	if f.isRunning {
		f.lockState.Unlock()
		return fmt.Errorf("already running")
	}

	var wg sync.WaitGroup
	wg.Add(3)

	msgCh := make(chan *Message)

	// close listener
	go func() {
		defer wg.Done()

		for {
			select {
			case <-f.endpoint.CloseCommandChannel():
				log.Info("internal stop")
				f.stopReading <- true
				f.stopProcessing <- true
				return
			case <-f.stopWaitClosing:
				log.Info("external stop")
				return
			}
		}
	}()

	// reader
	go func() {
		defer wg.Done()

		for {
			select {
			case <-f.stopReading:
				log.Info("external stop")
				return
			default:
				m, err := f.channel.Read()
				if err != nil {
					log.ErrorCtx("failed to read message", log.Context{"error": err})
					f.stopWaitClosing <- true
					f.stopProcessing <- true
					return
				}
				msgCh <- m
			}
		}
	}()

	// processor
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(f.pollingPeriod)

		for {
			select {
			case <-f.stopProcessing:
				log.Info("external stop")
				return
			case msg := <-msgCh:
				if err := f.endpoint.ProcessMessage(msg, f.channel.Write); err != nil {
					log.ErrorCtx("failed to process message", log.Context{"error": err})
					f.stopWaitClosing <- true
					f.stopReading <- true
					return
				}
			case t := <-ticker.C:
				if err := f.endpoint.Poll(t, f.channel.Write); err != nil {
					log.ErrorCtx("failed to poll", log.Context{"error": err})
					f.stopWaitClosing <- true
					f.stopReading <- true
					return
				}
			}
		}
	}()

	f.isRunning = true
	f.wgStop.Add(1)
	f.lockState.Unlock()
	wg.Wait()

	f.lockState.Lock()
	defer f.lockState.Unlock()

	f.isRunning = false

	f.channel.Close()
	f.isClosed = true

	f.wgStop.Done()
	return nil
}

// Close stops a full-duplex connection.
func (f *FullDuplex) Close() error {
	f.lockStop.Lock()
	defer f.lockStop.Unlock()

	f.lockState.Lock()
	isRunning := f.isRunning
	f.lockState.Unlock()

	if !isRunning {
		return fmt.Errorf("not running")
	}

	f.stopWaitClosing <- true
	f.stopReading <- true
	f.stopProcessing <- true

	f.wgStop.Wait()
	return nil
}

// IsRunning returns the running status of a full-duplex connection.
func (f *FullDuplex) IsRunning() bool {
	f.lockState.Lock()
	defer f.lockState.Unlock()

	return f.isRunning
}

// IsClosed returns the closed status of a full-duplex connection.
func (f *FullDuplex) IsClosed() bool {
	f.lockState.Lock()
	defer f.lockState.Unlock()

	return f.isClosed
}

// FullDuplexManager keeps track of a list of existing full-duplex connections.
type FullDuplexManager struct {
	connections []*FullDuplex
}

// NewFullDuplexManager creates a full-duplex connection list manager.
func NewFullDuplexManager() *FullDuplexManager {
	return new(FullDuplexManager)
}

// AddConnection appends a connection to the list of managed connections
// and runs it in a parallel goroutine.
func (m *FullDuplexManager) AddConnection(conn *FullDuplex) {
	m.connections = append(m.connections, conn)
	go conn.Run()
}

// CloseConnections closes all managed connections.
func (m *FullDuplexManager) CloseConnections() {
	for _, c := range m.connections {
		err := c.Close()
		if err != nil {
			log.ErrorCtx("manager - failed to close connection", log.Context{"error": err})
		}
	}
}
