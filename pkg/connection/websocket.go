package connection

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"

	"github.com/popescu-af/saas-y/pkg/log"
)

func newWebSocketChannel(c *websocket.Conn) *Channel {
	return &Channel{
		Read: func() (*Message, error) {
			mt, message, err := c.ReadMessage()
			return &Message{Type: mt, Payload: message}, err
		},
		Write: func(m *Message) error {
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

// NewWebSocketClient creates a new websocket connection and a full-duplex
// connection controller on top of it, returning both. It is the responsibility
// of the caller to both close the websocket connection and to stop the controller
// in case it is runnning.
func NewWebSocketClient(url url.URL, handler FullDuplexEndpoint, pollingPeriod time.Duration) (*FullDuplex, func(), error) {
	c, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		log.ErrorCtx("dial", log.Context{"error": err})
		return nil, nil, err
	}

	channel := newWebSocketChannel(c)
	conn := NewFullDuplex(handler, channel, pollingPeriod)
	return conn, func() { c.Close() }, nil
}

var upgrader = websocket.Upgrader{}

// NewWebSocketServer does the same as NewWebSocketClient, but from a server point of view.
// It creates the websocket by upgrading the HTTP request, compared to the client version,
// which creates the websocket by dialing an HTTP endpoint accepting the protocol.
func NewWebSocketServer(w http.ResponseWriter, r *http.Request, handler FullDuplexEndpoint, pollingPeriod time.Duration) (*FullDuplex, func(), error) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.ErrorCtx("upgrade", log.Context{"error": err})
		return nil, nil, err
	}

	channel := newWebSocketChannel(c)
	conn := NewFullDuplex(handler, channel, pollingPeriod)
	return conn, func() { c.Close() }, nil
}
