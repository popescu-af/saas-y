package connection

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type TestEndpoint struct {
}

func (t *TestEndpoint) ProcessMessage(*Message, WriteFn) error {
	return nil
}

func (t *TestEndpoint) Poll(time.Time, WriteFn) error {
	return nil
}

func TestLifetime(t *testing.T) {
	endp := &TestEndpoint{}
	channelIsOpen := true
	ch := Channel{
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

	conn := NewFullDuplex(endp, ch, time.Second)

	go func() {
		err := conn.Run()
		require.NoError(t, err, "unexpected error returned by run")
	}()
	time.Sleep(time.Duration(100) * time.Millisecond)

	require.True(t, conn.IsRunning(), "unexpected not running connection")

	err := conn.Stop()
	require.NoError(t, err, "unexpected error when stopping")
	require.False(t, conn.IsRunning(), "unexpected running connection")

	err = conn.Stop()
	require.Error(t, err, "unexpected successful stop when connection wasn't running")
}
