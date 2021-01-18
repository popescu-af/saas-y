package talk

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Messenger is an abstraction for multiple client - single service communications,
// for the purpose of testing the service.
type Messenger struct {
	Chats []*ChatData
}

// NewMessenger spawns n given communication channels between a service and n clients.
func NewMessenger(t *testing.T, n int, agentFuncsCreator AgentFuncsCreator, clientName, serviceName string) *Messenger {
	result := new(Messenger)

	for i := 0; i < n; i++ {
		chatData := new(ChatData)

		// Instantiate service stuff
		var recvService ReceiveFun
		var err error
		recvService, chatData.releaseService, err = agentFuncsCreator(t)
		require.NoError(t, err, "unexpected failure when creating service listener")

		// Instantiate client stuff
		recvClient := func(m *Message, reply SendFun) {
			chatData.ClientInbox = m
		}

		chatData.Conns, chatData.Agents, chatData.closeChat = NewChatSetup(recvService, recvClient)
		result.Chats = append(result.Chats, chatData)
	}

	return result
}

// Chat is a helper struct managing the communication between two talking agents.
type Chat struct {
	wg     sync.WaitGroup
	agent1 *Agent
	agent2 *Agent
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

// NewFakeChannel creates a new fake channel instance.
func NewFakeChannel() *FakeChannel {
	result := &FakeChannel{
		ch: make([]chan *Message, 2),
	}
	result.ch[0] = make(chan *Message)
	result.ch[1] = make(chan *Message)
	return result
}

// NewFakeConnection returns a new fake connection instance.
func NewFakeConnection(channel *FakeChannel, index int) *FakeConnection {
	return &FakeConnection{
		channel: channel,
		index:   index,
	}
}

// AgentFuncsCreator creates channel listeners for the service side, for testing.
type AgentFuncsCreator func(t *testing.T) (ReceiveFun, func(), error)

// Stop closes all chats in a messenger and releases the service.
func (m *Messenger) Stop() {
	for _, c := range m.Chats {
		c.closeChat()
		c.releaseService()
	}
}

// Wait waits for the chat to be over.
func (c *Chat) Wait() {
	c.wg.Wait()
}

// ChatData holds information about a chat.
type ChatData struct {
	Conns          []*FakeConnection
	Agents         []*Agent
	ClientInbox    *Message
	closeChat      func()
	releaseService func()
}

// ExpectServiceMessage sets the expectation for receiving a message on servce-side.
func (c *ChatData) ExpectServiceMessage() {
	c.Conns[0].ExpectMessages(1)
}

// WaitForServiceMessages waits for all the message expectations on service-side to be fulfilled.
func (c *ChatData) WaitForServiceMessages() {
	c.Conns[0].WaitForMessages()
}

// ExpectClientMessage sets the expectation for receiving a message on client-side.
func (c *ChatData) ExpectClientMessage() {
	c.Conns[1].ExpectMessages(1)
}

// WaitForClientMessages waits for all the message expectations on client-side to be fulfilled.
func (c *ChatData) WaitForClientMessages() {
	c.Conns[1].WaitForMessages()
}

// FakeChannel represents a theoretical communication channel.
// between two parties. It is useful for mocking real communication channels.
// In a FakeChannel, there are two Message channels. One party reads from the first
// channel and writes to the second channel. The other party does the opposite.
type FakeChannel struct {
	ch []chan *Message
}

// FakeConnection represents a theoretical communication endpoint.
// It is useful for mocking real communication endpoints.
type FakeConnection struct {
	channel        *FakeChannel
	index          int
	messageArrived sync.WaitGroup
	mtxCallback    sync.Mutex
	callback       func(*Message, bool)
}

func (c *FakeConnection) cb(m *Message, closed bool) {
	c.mtxCallback.Lock()
	defer c.mtxCallback.Unlock()
	if c.callback != nil {
		c.callback(m, closed)
	}
}

// Subscribe subscribes the given callback to the connection's
// message arrival event loop.
func (c *FakeConnection) Subscribe(cb func(*Message, bool)) {
	c.mtxCallback.Lock()
	defer c.mtxCallback.Unlock()
	c.callback = cb
}

// Write writes a message on the fake channel.
func (c *FakeConnection) Write(m *Message) error {
	// async posting of the message
	go func() {
		c.channel.ch[1-c.index] <- m
	}()
	return nil
}

// Close closes the fake connection by writing a close message on the fake channel.
func (c *FakeConnection) Close() {
	c.Write(&Message{Type: CloseMessage})
}

// ExpectMessages expectats for 'count' messages.
func (c *FakeConnection) ExpectMessages(count int) {
	c.messageArrived.Add(count)
	go func() {
		for i := 0; i < count; i++ {
			m := <-c.channel.ch[c.index]
			c.cb(m, m.Type == CloseMessage)
			c.messageArrived.Done()
		}
	}()
}

// WaitForMessages waits for the expected messages.
func (c *FakeConnection) WaitForMessages() {
	c.messageArrived.Wait()
}
