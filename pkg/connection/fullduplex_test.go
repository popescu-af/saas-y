package connection

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var channelIsOpen bool
var defaultChannel = &Channel{
	Read: func() (*Message, error) {
		if !channelIsOpen {
			return nil, fmt.Errorf("error reading from closed channel")
		}
		return &Message{Type: 0, Payload: []byte("dummy message")}, nil
	},
	Write: func(m *Message) error {
		// do nothing
		return nil
	},
	Close: func() {
		channelIsOpen = false
	},
}

type TestEndpoint struct {
	channelClosed                 bool
	closeCh                       chan bool
	stopInternallyFromProcessFunc bool
	stopInternallyFromPollFunc    bool
}

func (t *TestEndpoint) ProcessMessage(*Message, WriteFn) error {
	if !t.channelClosed && t.stopInternallyFromProcessFunc {
		t.closeCh <- true
		t.channelClosed = true
		close(t.closeCh)
	}
	return nil
}

func (t *TestEndpoint) Poll(time.Time, WriteFn) error {
	if !t.channelClosed && t.stopInternallyFromPollFunc {
		t.closeCh <- true
		t.channelClosed = true
		close(t.closeCh)
	}
	return nil
}

func (t *TestEndpoint) CloseCommandChannel() chan bool {
	return t.closeCh
}

func NewTestEndpoint(stopInternallyFromProcessFunc, stopInternallyFromPollFunc bool) *TestEndpoint {
	return &TestEndpoint{
		closeCh:                       make(chan bool),
		stopInternallyFromProcessFunc: stopInternallyFromProcessFunc,
		stopInternallyFromPollFunc:    stopInternallyFromPollFunc,
	}
}

func TestLifetime(t *testing.T) {
	endp := NewTestEndpoint(false, false)

	channelIsOpen = true
	conn := NewFullDuplex(endp, defaultChannel, time.Second)

	go func() {
		err := conn.Run()
		require.NoError(t, err, "unexpected error returned by run")
	}()
	time.Sleep(time.Duration(100) * time.Millisecond)

	require.True(t, conn.IsRunning(), "unexpected not running connection")
	require.False(t, conn.IsClosed(), "unexpected closed connection")

	err := conn.Close()
	require.NoError(t, err, "unexpected error when stopping")
	require.False(t, conn.IsRunning(), "unexpected running connection")

	require.False(t, conn.IsRunning(), "unexpected running connection")
	require.True(t, conn.IsClosed(), "unexpected not closed connection")

	err = conn.Close()
	require.Error(t, err, "unexpected successful close when connection wasn't running")

	err = conn.Run()
	require.Error(t, err, "unexpected successful run when connection should be closed")
}

func TestInternalStop(t *testing.T) {
	tests := [][]bool{
		{false, true},
		{true, false},
	}

	for _, tt := range tests {
		endp := NewTestEndpoint(tt[0], tt[1])

		channelIsOpen = true
		conn := NewFullDuplex(endp, defaultChannel, time.Second)

		err := conn.Run()
		require.NoError(t, err, "unexpected error returned by run")

		require.False(t, conn.IsRunning(), "unexpected running connection")
		require.True(t, conn.IsClosed(), "unexpected not closed connection")
	}
}
