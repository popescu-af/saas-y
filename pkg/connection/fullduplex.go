package connection

import (
	"fmt"
	"sync"

	"github.com/popescu-af/saas-y/pkg/log"
)

// Channel is the interface for a full duplex channel.
type Channel interface {
	Read() (*Message, error)
	Write(*Message) error
	Close()
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

// WriteOnChannelFunc can be called to write back on the channel
// from inside the processing message method.
type WriteOnChannelFunc func(*Message) error

// ErrorStop is the error a listener should return when processing a message,
// if the connection should be subsequently closed.
var ErrorStop = fmt.Errorf("listener: stop connection")

// ChannelListener is the interface for types that can process incoming messages from
// a full duplex channel. The listener can react to messages with a reply or by closing the channel.
type ChannelListener interface {
	// ProcessMessage should process the given message and react accordingly,
	// by changing state and/or writing something back or none of them.
	ProcessMessage(*Message, WriteOnChannelFunc) error
}

// FullDuplex is a full-duplex connection that takes a channel listener
// and a channel. It handles messages arriving on the channel through the listener
// and also handles sending messages and closing the communication.
type FullDuplex struct {
	listener  ChannelListener
	channel   Channel
	stopCh    chan bool
	wgStop    sync.WaitGroup
	lockState sync.Mutex
	lockStop  sync.Mutex
	isRunning bool
	isClosed  bool
}

// NewFullDuplex creates a new, inactive full-duplex connection.
// Call Run to run it.
func NewFullDuplex(listener ChannelListener, channel Channel) *FullDuplex {
	return &FullDuplex{
		listener: listener,
		channel:  channel,
		stopCh:   make(chan bool),
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
	wg.Add(2)

	// processor <- reader
	msgCh := make(chan *Message)

	// reader
	go func() {
		defer wg.Done()
		defer log.InfoCtx("reading done", map[string]interface{}{"instance": fmt.Sprintf("%p", f)})

		for {
			msg, err := f.channel.Read()
			if err != nil {
				log.ErrorCtx("failed to read message", log.Context{"error": err})
				return
			}

			msgCh <- msg

			if msg.Type == CloseMessage {
				return
			}
		}
	}()

	// processor
	go func() {
		defer wg.Done()
		defer log.InfoCtx("processing done", map[string]interface{}{"instance": fmt.Sprintf("%p", f)})

		writeOnChannel := func(m *Message) error {
			return f.channel.Write(m)
		}

		for {
			select {
			case <-f.stopCh:
				log.Info("external stop")
				f.channel.Close()
				return
			case msg := <-msgCh:
				switch msg.Type {
				case CloseMessage:
					log.Info("channel closed by the other party")
					return
				case PingMessage:
					if err := writeOnChannel(&Message{Type: PongMessage}); err != nil {
						log.ErrorCtx("failed to send pong", log.Context{"error": err})
						return
					}
				case PongMessage:
					// do nothing for now
				default:
					if err := f.listener.ProcessMessage(msg, writeOnChannel); err != nil {
						log.ErrorCtx("failed to process message", log.Context{"error": err})
						return
					}
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
	f.isClosed = true

	f.wgStop.Done()
	return nil
}

// SendMessage sends a message on the full duplex channel.
func (f *FullDuplex) SendMessage(m *Message) error {
	return f.channel.Write(m)
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

	f.stopCh <- true

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
