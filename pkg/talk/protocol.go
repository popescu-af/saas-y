package talk

import (
	"encoding/json"
	"fmt"
)

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

// Message is the unit of transmission.
type Message struct {
	Type    int
	Payload []byte
}

// SimpleMessage creates a simple text message.
func SimpleMessage(text string) *Message {
	return &Message{
		Type:    TextMessage,
		Payload: []byte(text),
	}
}

// ToTextMessage converts a JSON-annotated struct to a Message of type Text.
// It returns a nil value if the marshaling does not succeed.
func ToTextMessage(v interface{}) (*Message, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return &Message{
		Type:    TextMessage,
		Payload: b,
	}, nil
}

// FromTextMessage creates an instance of a JSON-annotated type from a
// Message of type Text.
func FromTextMessage(m *Message, v interface{}) error {
	if m == nil {
		return fmt.Errorf("empty message")
	}
	return json.Unmarshal(m.Payload, v)
}

// Connection is the interface for an endpoint.
type Connection interface {
	Subscribe(func(*Message))
	Write(*Message) error
	Close()
}

// SendFun is the signature of the function that sends a message.
type SendFun func(*Message)

// ReceiveFun is the signature of the function that handles a received message.
type ReceiveFun func(*Message, SendFun)
