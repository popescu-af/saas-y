package talk

import (
	"sync"
)

// FakeChannel represents a theoretical communication channel.
// between two parties. It is useful for mocking real communication channels.
// In a FakeChannel, there are two Message channels. One party reads from the first
// channel and writes to the second channel. The other party does the opposite.
type FakeChannel struct {
	ch []chan *Message
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

// NewFakeConnection returns a new fake connection instance.
func NewFakeConnection(channel *FakeChannel, index int) *FakeConnection {
	return &FakeConnection{
		channel: channel,
		index:   index,
	}
}
