package connection

import (
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
			for channelIsOpen {
				// block until channelIsOpen == false
			}
			return nil, nil
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

	go conn.Run()
	time.Sleep(time.Duration(100) * time.Millisecond)

	require.True(t, conn.IsRunning(), "unexpected not running connection")

	conn.Stop()
	require.False(t, conn.IsRunning(), "unexpected running connection")
}
