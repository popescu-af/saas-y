package connection

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testChannel struct {
	mock.Mock
	failReading        bool
	failWriting        bool
	messageType        int
	channelIsClosed    bool
	lastWrittenMessage *Message
}

var errRead = fmt.Errorf("error reading message")

func (t *testChannel) Read() (*Message, error) {
	if t.channelIsClosed {
		return nil, fmt.Errorf("error reading from closed channel")
	} else if t.failReading {
		t.Called()
		return nil, errRead
	}
	return &Message{Type: t.messageType, Payload: []byte("dummy message")}, nil
}

var errWrite = fmt.Errorf("error writing message")

func (t *testChannel) Write(m *Message) error {
	t.lastWrittenMessage = m
	t.Called()

	if t.failWriting {
		return errWrite
	}
	return nil
}

func (t *testChannel) Close() {
	t.channelIsClosed = true
}

type testListener struct {
	mock.Mock
	stopInternally            bool
	replyToProcessMessageOnce bool
}

func (t *testListener) ProcessMessage(m *Message, write WriteOnChannelFunc) error {
	if t.stopInternally {
		t.Called()
		return ErrorStop
	} else if t.replyToProcessMessageOnce {
		t.replyToProcessMessageOnce = false
		return write(nil)
	}
	return nil
}

func setupTest(stopInternally, failReading, failWriting, replyToProcessMessageOnce bool, msgType int) (*FullDuplex, *testListener, *testChannel) {
	listener := &testListener{
		stopInternally:            stopInternally,
		replyToProcessMessageOnce: replyToProcessMessageOnce,
	}
	channel := &testChannel{
		failReading: failReading,
		failWriting: failWriting,
		messageType: msgType,
	}
	return NewFullDuplex(listener, channel), listener, channel
}

func TestLifetime(t *testing.T) {
	conn, _, channel := setupTest(false, false, false, true, TextMessage)

	// it will happen as a single reply to the first processed message
	channel.On("Write")

	go func() {
		err := conn.Run()
		require.NoError(t, err, "unexpected error returned by run")
	}()
	time.Sleep(time.Duration(100) * time.Millisecond)

	require.True(t, conn.IsRunning(), "unexpected not running connection")
	require.False(t, conn.IsClosed(), "unexpected closed connection")

	// see comment above
	channel.AssertCalled(t, "Write")

	err := conn.Run()
	require.Error(t, err, "unexpected successful second run")

	channel.On("Write")
	conn.SendMessage(nil)
	channel.AssertCalled(t, "Write")

	err = conn.Close()
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
	conn, listener, _ := setupTest(true, false, false, false, TextMessage)

	listener.On("ProcessMessage").Return(ErrorStop)
	conn.Run()
	listener.AssertCalled(t, "ProcessMessage")
}

func TestReadError(t *testing.T) {
	conn, _, channel := setupTest(false, true, false, false, TextMessage)

	channel.On("Read").Return(errRead)
	conn.Run()
	channel.AssertCalled(t, "Read")
}

func TestPingReply(t *testing.T) {
	conn, _, channel := setupTest(false, false, false, false, PingMessage)

	channel.On("Write")
	go func() {
		err := conn.Run()
		require.NoError(t, err, "unexpected error returned by run")
	}()
	time.Sleep(time.Duration(100) * time.Millisecond)

	channel.AssertCalled(t, "Write")
	conn.Close()

	require.Equal(t, PongMessage, channel.lastWrittenMessage.Type, "unexpected last sent message type: %d", channel.lastWrittenMessage.Type)
}

func TestPingReplyFailed(t *testing.T) {
	conn, _, channel := setupTest(false, false, true, false, PingMessage)

	channel.On("Write").Return(errWrite)
	conn.Run()
	channel.AssertCalled(t, "Write")

	require.Equal(t, PongMessage, channel.lastWrittenMessage.Type, "unexpected last sent message type: %d", channel.lastWrittenMessage.Type)
}

func TestCloseMessage(t *testing.T) {
	conn, _, _ := setupTest(false, false, false, false, CloseMessage)

	conn.Run()

	require.False(t, conn.IsRunning(), "unexpected running connection")
	require.True(t, conn.IsClosed(), "unexpected open connection")
}
