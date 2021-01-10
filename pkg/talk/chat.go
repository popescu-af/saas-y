package talk

import (
	"sync"
	"time"
)

// Chat is a helper struct managing two talking agents.
type Chat struct {
	wg     sync.WaitGroup
	agent1 *Agent
	agent2 *Agent
}

// Wait waits for the chat to be over.
func (c *Chat) Wait() {
	c.wg.Wait()
}

// NewChat creates a new chat between two talking agents.
func NewChat(agent1, agent2 *Agent) *Chat {
	c := &Chat{
		agent1: agent1,
		agent2: agent2,
	}

	c.wg.Add(2)
	go func() {
		defer c.wg.Done()
		c.agent1.Listen()
	}()
	go func() {
		defer c.wg.Done()
		c.agent2.Listen()
	}()

	// wait for the agents to start
	time.Sleep(time.Duration(10) * time.Millisecond)
	return c
}

// NewChatSetup creates the whole setup required for a chat between two talking agents.
func NewChatSetup(recv1, recv2 ReceiveFun) ([]*FakeConnection, []*Agent, func()) {
	channel := NewFakeChannel()
	conns := []*FakeConnection{
		NewFakeConnection(channel, 0),
		NewFakeConnection(channel, 1),
	}
	agents := []*Agent{
		NewAgent(conns[0], recv1),
		NewAgent(conns[1], recv2),
	}

	chat := NewChat(agents[0], agents[1])
	return conns, agents, func() {
		agents[0].Stop()
		agents[1].Stop()
		chat.Wait()
	}
}
