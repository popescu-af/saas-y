package connection

import (
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"

	"github.com/popescu-af/saas-y/pkg/log"
)

type webSocketChannel struct {
	wsConn *websocket.Conn
}

func (w *webSocketChannel) Read() (*Message, error) {
	mt, message, err := w.wsConn.ReadMessage()
	return &Message{Type: mt, Payload: message}, err
}

func (w *webSocketChannel) Write(m *Message) error {
	return w.wsConn.WriteMessage(m.Type, m.Payload)
}

func (w *webSocketChannel) Close() {
	closeMsgFmt := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	err := w.wsConn.WriteMessage(websocket.CloseMessage, closeMsgFmt)
	if err != nil {
		log.ErrorCtx("write close message", log.Context{"error": err})
	}
	w.wsConn.Close()
}

func newWebSocketChannel(c *websocket.Conn) Channel {
	return &webSocketChannel{wsConn: c}
}

// NewWebSocketClient creates a new websocket connection and a full-duplex
// connection on top of it.
func NewWebSocketClient(url url.URL, listener ChannelListener) (*FullDuplex, error) {
	c, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		log.ErrorCtx("dial", log.Context{"error": err})
		return nil, err
	}

	channel := newWebSocketChannel(c)
	conn := NewFullDuplex(listener, channel, "client")
	return conn, nil
}

// TODO: Properly instantiate the upgrader, with the possibility of
// customizing the origin checker function
var upgrader = func() websocket.Upgrader {
	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	return upgrader
}()

// NewWebSocketServer does the same as NewWebSocketClient, but from a server point of view.
func NewWebSocketServer(w http.ResponseWriter, r *http.Request, listener ChannelListener) (*FullDuplex, error) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.ErrorCtx("upgrade", log.Context{"error": err})
		return nil, err
	}

	channel := newWebSocketChannel(c)
	conn := NewFullDuplex(listener, channel, "server")
	return conn, nil
}
