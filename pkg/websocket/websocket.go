package websocket

import (
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/mailru/easygo/netpoll"

	"github.com/popescu-af/saas-y/pkg/log"
	"github.com/popescu-af/saas-y/pkg/talk"
)

// connection is a wrapper over a web socket connection
// that implements the talk.Connection interface.
type connection struct {
	conn   *websocket.Conn
	poller netpoll.Poller
}

// Subscribe subscribes the given callback to the poller read events.
func (c *connection) Subscribe(cb func(*talk.Message)) {
	fd := netpoll.Must(netpoll.HandleRead(c.conn.UnderlyingConn()))
	c.poller.Start(fd, func(ev netpoll.Event) {
		if ev&netpoll.EventReadHup != 0 {
			c.poller.Stop(fd)
			return
		}

		mt, message, err := c.conn.ReadMessage()
		if err != nil {
			log.ErrorCtx("read message", log.Context{"error": err})
			return
		}

		cb(&talk.Message{Type: mt, Payload: message})
	})
}

// Write sends a message over the web socket connection.
func (c *connection) Write(m *talk.Message) error {
	return c.conn.WriteMessage(m.Type, m.Payload)
}

// Close closes the web socket connection.
func (c *connection) Close() {
	c.conn.Close()
}

func newConnection(c *websocket.Conn) (*connection, error) {
	poller, err := netpoll.New(nil)
	if err != nil {
		return nil, err
	}

	return &connection{
		conn:   c,
		poller: poller,
	}, nil
}

// NewClient creates a talking agent wrapping a websocket client.
func NewClient(url url.URL, recv talk.ReceiveFun) (*talk.Agent, error) {
	c, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		log.ErrorCtx("dial", log.Context{"error": err})
		return nil, err
	}

	conn, err := newConnection(c)
	if err != nil {
		return nil, err
	}

	return talk.NewAgent(conn, recv), nil
}

// TODO: Properly instantiate the upgrader, with the possibility of
// customizing the origin checker function
var upgrader = func() websocket.Upgrader {
	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	return upgrader
}()

// NewServer creates a talking agent wrapping a websocket server.
func NewServer(w http.ResponseWriter, r *http.Request, recv talk.ReceiveFun) (*talk.Agent, error) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.ErrorCtx("upgrade", log.Context{"error": err})
		return nil, err
	}

	conn, err := newConnection(c)
	if err != nil {
		return nil, err
	}

	return talk.NewAgent(conn, recv), nil
}
