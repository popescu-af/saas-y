package websocket

import (
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/mailru/easygo/netpoll"

	"github.com/popescu-af/saas-y/pkg/log"
	"github.com/popescu-af/saas-y/pkg/talk"
)

type pollerPool struct {
	pollers []netpoll.Poller
	index   int
	count   int
}

func (p *pollerPool) Next() netpoll.Poller {
	p.index = (p.index + 1) % p.count
	return p.pollers[p.index]
}

var defaultPollerPool = initPollerPool(32)

// ResizePollerPool initializes the poller pool to the desired number
// of netpoll pollers, which will be used by the websocket connections.
func ResizePollerPool(n int) {
	defaultPollerPool = initPollerPool(n)
}

func initPollerPool(n int) *pollerPool {
	if n < 1 {
		n = 1
	}

	p := &pollerPool{
		pollers: make([]netpoll.Poller, n),
		index:   n - 1,
		count:   n,
	}

	var err error
	for i := 0; i < n; i++ {
		p.pollers[i], err = netpoll.New(nil)
		if err != nil {
			log.ErrorCtx("cannot create poller", log.Context{"error": err})
			return nil
		}
	}
	return p
}

// connection is a wrapper over a web socket connection
// that implements the talk.Connection interface.
type connection struct {
	conn *websocket.Conn
}

// Subscribe subscribes the given callback to the poller read events.
func (c *connection) Subscribe(cb func(*talk.Message, bool)) {
	fd := netpoll.Must(netpoll.HandleRead(c.conn.UnderlyingConn()))
	poller := defaultPollerPool.Next()
	poller.Start(fd, func(ev netpoll.Event) {
		unsubscribe := false
		defer func() {
			if unsubscribe {
				cb(nil, true)
				poller.Stop(fd)
			}
		}()

		if ev&netpoll.EventReadHup != 0 {
			unsubscribe = true
			return
		}

		mt, message, err := c.conn.ReadMessage()
		if err != nil {
			log.ErrorCtx("read message", log.Context{"error": err})
			unsubscribe = true
			return
		}

		cb(&talk.Message{Type: mt, Payload: message}, false)
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
	return &connection{conn: c}, nil
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
