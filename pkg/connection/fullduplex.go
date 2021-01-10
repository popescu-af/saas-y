package connection

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/mailru/easygo/netpoll"
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

// Full duplex protocol message codes.
const (
	// Success is used for replying when the previous request was processed successfully,
	// without any particular result being sent back. When results are involved, positive
	// codes should be defined by implementations, depending on their use case.
	Success = 0
	// InvalidMessage is received when the last request contained a malformed message,
	// from the point of view of the service implementation.
	InvalidMessage = -1
	// KeyCollision is used when a key already exists and cannot be acted upon.
	KeyCollision = -2
	// NotFound is used when the specified resource is non-existent.
	NotFound = -3
	// Unauthorized is used when requested access to the specified resource needs authorization.
	Unauthorized = -4
	// NotAllowed is used when requested access to the specified resource is not allowed.
	NotAllowed = -5
	// Timeout is used when a request times out.
	Timeout = -6
	// InternalError is used when something goes wrong in the implementation.
	InternalError = -7
)

// WriteOnChannelFunc can be called to write back on the channel
// from inside the processing message method.
type WriteOnChannelFunc func(*Message)

// ErrorStop is the error a listener should return when processing a message,
// if the connection should be subsequently closed.
var ErrorStop = fmt.Errorf("listener: stop connection")

// ChannelListener is the interface for types that can process incoming messages from
// a full duplex channel. The listener can react to messages with a reply or by closing the channel.
type ChannelListener interface {
	// ProcessMessage should process the given message and react accordingly,
	// by changing state and/or writing something back or none of them.
	ProcessMessage(*Message, WriteOnChannelFunc)
}

// FullDuplex is a full-duplex connection that takes a channel listener
// and a channel. It handles messages arriving on the channel through the listener
// and also handles sending messages and closing the communication.
type FullDuplex struct {
	name     string
	listener ChannelListener
	channel  Channel
	writeCh  chan *Message
	stopCh   chan bool
	finished sync.WaitGroup
}

// NewFullDuplex creates a new, inactive full-duplex connection.
// Call Run to run it.
func NewFullDuplex(listener ChannelListener, channel Channel, name string) *FullDuplex {
	return &FullDuplex{
		name:     name,
		listener: listener,
		channel:  channel,
		writeCh:  make(chan *Message),
		stopCh:   make(chan bool),
	}
}

// Run is a blocking function that handles messages arriving on the channel.
func (f *FullDuplex) Run() error {
	defer log.DebugCtx("channel closed", log.Context{"name": f.name})
	defer f.finished.Done()

	f.finished.Add(1)

	poller, err := netpoll.New(nil)
	if err != nil {
		return err
	}

	// TODO: make WS-agnostic
	fd := netpoll.Must(netpoll.HandleRead(f.channel.(*webSocketChannel).wsConn.UnderlyingConn()))

	readCh := make(chan *Message)
	poller.Start(fd, func(ev netpoll.Event) {
		if ev&netpoll.EventReadHup != 0 {
			poller.Stop(fd)
			f.stopCh <- true
			return
		}

		msg, err := f.channel.Read()
		if err != nil {
			log.ErrorCtx("failed to read message", log.Context{"name": f.name, "error": err})
			return
		}
		readCh <- msg
	})

	handleRead := func(msg *Message) {
		switch msg.Type {
		case CloseMessage:
			log.DebugCtx("channel closed by the other party", log.Context{"name": f.name})
			f.stopCh <- true
		case PingMessage:
			log.InfoCtx("received ping", log.Context{"name": f.name})
			f.writeCh <- &Message{Type: PongMessage}
		case PongMessage:
			log.InfoCtx("received pong", log.Context{"name": f.name})
		default:
			log.DebugCtx("processing with payload", log.Context{"name": f.name, "msg_payload": string(msg.Payload)})
			f.listener.ProcessMessage(msg, func(msg *Message) {
				f.writeCh <- msg
			})
			log.DebugCtx("done processing", log.Context{"name": f.name})
		}
	}

	handleWrite := func(msg *Message) {
		log.DebugCtx("sending with payload", log.Context{"name": f.name, "msg_payload": string(msg.Payload)})
		if err := f.channel.Write(msg); err != nil {
			log.ErrorCtx("failed to send message", log.Context{"name": f.name, "error": err})
		}
	}

	// main loop
	for {
		select {
		case msg := <-readCh:
			go handleRead(msg)
		case msg := <-f.writeCh:
			handleWrite(msg)
		case <-f.stopCh:
			return nil
		}
	}
}

// SendMessage sends a message on the full duplex channel.
func (f *FullDuplex) SendMessage(m *Message) {
	f.writeCh <- m
}

// Close stops a full-duplex connection.
func (f *FullDuplex) Close() {
	defer f.finished.Wait()
	f.stopCh <- true
}

// // IsRunning returns the running status of a full-duplex connection.
// func (f *FullDuplex) IsRunning() bool {
// 	f.lockState.Lock()
// 	defer f.lockState.Unlock()

// 	return f.isRunning
// }

// // IsClosed returns the closed status of a full-duplex connection.
// func (f *FullDuplex) IsClosed() bool {
// 	f.lockState.Lock()
// 	defer f.lockState.Unlock()

// 	return f.isClosed
// }

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
		// err :=
		c.Close()
		// if err != nil {
		// 	log.ErrorCtx("manager - failed to close connection", log.Context{"error": err})
		// }
	}
}

// ToTextMessage converts a JSON-annotated struct to a Message of type Text.
// It returns a nil value if the marshaling does not succeed.
func ToTextMessage(v interface{}) *Message {
	b, err := json.Marshal(v)
	if err != nil {
		log.ErrorCtx("failed to marshal input", map[string]interface{}{"error": err})
		return nil
	}

	return &Message{
		Type:    TextMessage,
		Payload: b,
	}
}

// FromTextMessage creates an instance of a JSON-annotated type from a
// Message of type Text.
func FromTextMessage(m *Message, v interface{}) error {
	return json.Unmarshal(m.Payload, v)
}
