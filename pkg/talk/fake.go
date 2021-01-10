package talk

import (
	"sync"
	"time"
)

// FakeChannel represents a theoretical communication channel.
// It is useful for mocking real communication channels.
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
	exitCh         chan struct{}
	exited         sync.WaitGroup
}

// Subscribe subscribes the given callback to the connection's
// message arrival event loop.
func (c *FakeConnection) Subscribe(cb func(*Message)) {
	c.exited.Add(1)

	go func() {
		defer c.exited.Done()
		for {
			select {
			case m := <-c.channel.ch[c.index]:
				cb(m)
				c.messageArrived.Done()
			case <-c.exitCh:
				return
			}
		}
	}()

	// wait for select to be in place
	time.Sleep(time.Duration(10) * time.Millisecond)
}

// Write writes a message on the fake channel.
func (c *FakeConnection) Write(m *Message) error {
	c.channel.ch[1-c.index] <- m
	return nil
}

// Close closes the fake connection.
func (c *FakeConnection) Close() {
	defer c.exited.Wait()
	c.exitCh <- struct{}{}
}

// ExpectMessages sets expectation for 'count' messages.
func (c *FakeConnection) ExpectMessages(count int) {
	c.messageArrived.Add(count)
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
		exitCh:  make(chan struct{}),
	}
}
