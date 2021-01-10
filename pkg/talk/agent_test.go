package talk

import (
	"testing"
)

func TestSimpleMessageArrival(t *testing.T) {
	recv1 := func(m *Message, reply SendFun) {
	}
	recv2 := func(m *Message, reply SendFun) {
	}
	conns, agents, stopChat := NewChatSetup(recv1, recv2)
	defer stopChat()

	conns[1].ExpectMessages(1)
	agents[0].Say(SimpleMessage("hello, A1!"))
	conns[1].WaitForMessages()
}

func TestStopAgentMultipleTimes(t *testing.T) {
	recv1 := func(m *Message, reply SendFun) {
	}
	recv2 := func(m *Message, reply SendFun) {
	}
	_, agents, stopChat := NewChatSetup(recv1, recv2)
	defer stopChat()

	agents[0].Stop()
	agents[0].Stop()
}
