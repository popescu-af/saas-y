package talk

import (
	"sync"

	"github.com/popescu-af/saas-y/pkg/worker"
)

// WorkerPool is a worker pool, which if initialised, is used
// instead of launching new goroutines every time when needed.
var WorkerPool *worker.Pool

// Agent is the class for talking agents.
// A talking agent both accepts and sends messages.
type Agent struct {
	mtxConn   sync.Mutex
	conn      Connection
	recv      ReceiveFun
	exitCh    chan struct{}
	exited    sync.WaitGroup
	mtxStatus sync.Mutex
	running   bool
}

func (a *Agent) say(m *Message) {
	a.mtxConn.Lock()
	defer a.mtxConn.Unlock()
	a.conn.Write(m)
	if m.Type == CloseMessage {
		a.conn.Close()
	}
}

func (a *Agent) setRunning(r bool) {
	a.mtxStatus.Lock()
	defer a.mtxStatus.Unlock()
	a.running = r
}

func (a *Agent) stop() bool {
	if !a.IsRunning() {
		return false
	}
	a.exitCh <- struct{}{}
	a.exited.Wait()
	return true
}

// IsRunning returns true if the agent is still running.
// An stopped agent is not usable anymore.
func (a *Agent) IsRunning() bool {
	a.mtxStatus.Lock()
	defer a.mtxStatus.Unlock()
	return a.running
}

// Listen blocks the current goroutine and makes the agent listen
// for incoming messages until Stop() is called.
func (a *Agent) Listen() {
	defer a.exited.Done()

	a.exited.Add(1)
	a.setRunning(true)
	defer a.setRunning(false)

	a.conn.Subscribe(func(m *Message, closed bool) {
		if closed {
			a.stop()
			return
		}
		a.recv(m, a.Say)
	})

	for range a.exitCh {
		break
	}
}

// Say sends a message to the peer agent.
func (a *Agent) Say(m *Message) {
	f := func() {
		a.say(m)
	}
	if WorkerPool != nil {
		WorkerPool.Post(f)
	} else {
		go f()
	}
}

// Stop halts listening and closes the connection.
func (a *Agent) Stop() {
	if a.stop() {
		a.say(&Message{Type: CloseMessage})
	}
}

// NewAgent creates a new talking agent.
// The agent will take care of closing the connection when stopped.
func NewAgent(c Connection, r ReceiveFun) *Agent {
	return &Agent{
		conn:   c,
		recv:   r,
		exitCh: make(chan struct{}),
	}
}
