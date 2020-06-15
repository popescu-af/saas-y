package connection

import (
	"github.com/gorilla/websocket"

	"github.com/popescu-af/saas-y/pkg/log"
)

// NewWebSocketChannel returns an instance of websocket channel.
func NewWebSocketChannel(c *websocket.Conn) *Channel {
	return &connection.Channel{
		Read: func() (*connection.Message, error) {
			mt, message, err := c.ReadMessage()
			return &connection.Message{Type: mt, Payload: message}, err
		},
		Write: func(m *connection.Message) error {
			return c.WriteMessage(m.Type, m.Payload)
		},
		Close: func() {
			closeMsgFmt := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
			err := c.WriteMessage(websocket.CloseMessage, closeMsgFmt)
			if err != nil {
				log.ErrorCtx("write close message", log.Context{"error": err})
			}
		},
	}
}
