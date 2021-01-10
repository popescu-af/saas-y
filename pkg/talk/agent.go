package talk

import "sync"

// Agent is the class for talking agents.
type Agent struct {
	conn    Connection
	recv    ReceiveFun
	msgCh   chan *Message
	exitCh  chan struct{}
	exited  sync.WaitGroup
	running bool
}

// Listen listens for incoming messages until Stop() is called.
func (a *Agent) Listen() {
	defer a.exited.Done()

	a.exited.Add(1)
	a.running = true

	a.conn.Subscribe(func(m *Message) {
		a.recv(m, a.Say)
	})

	for {
		select {
		case m := <-a.msgCh:
			a.conn.Write(m)
		case <-a.exitCh:
			a.running = false
			return
		}
	}
}

// Say sends a message.
func (a *Agent) Say(m *Message) {
	go func() {
		a.msgCh <- m
	}()
}

// Stop closes the connection and stops listening.
func (a *Agent) Stop() {
	if !a.running {
		return
	}

	defer a.exited.Wait()
	a.conn.Close()
	a.exitCh <- struct{}{}
}

// NewAgent creates a new talking agent.
// The agent will take care of closing the connection when stopped.
func NewAgent(c Connection, r ReceiveFun) *Agent {
	return &Agent{
		conn:   c,
		recv:   r,
		msgCh:  make(chan *Message),
		exitCh: make(chan struct{}),
	}
}
