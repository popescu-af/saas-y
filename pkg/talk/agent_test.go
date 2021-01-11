package talk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimpleMessageArrival(t *testing.T) {
	recv0 := func(m *Message, reply SendFun) {
	}
	recv1 := func(m *Message, reply SendFun) {
	}
	conns, agents, stopChat := NewChatSetup(recv0, recv1)
	defer stopChat()

	conns[1].ExpectMessages(1)
	agents[0].Say(SimpleMessage("hello, A1!"))
	conns[1].WaitForMessages()
}

func TestStopAgentMultipleTimes(t *testing.T) {
	recv0 := func(m *Message, reply SendFun) {
	}
	recv1 := func(m *Message, reply SendFun) {
	}
	_, agents, stopChat := NewChatSetup(recv0, recv1)
	defer stopChat()

	agents[0].Stop()
	agents[0].Stop()
}

func TestSimpleTalk(t *testing.T) {
	a0InitMessage := "A1, how are you?"
	a1Reply := "Here is A1, I am well! How are you, A0?"
	a0Reply := "I'm well too!"
	a1ReceivedIndex := 0

	recv0 := func(m *Message, reply SendFun) {
		require.Equal(t, a1Reply, string(m.Payload))
	}
	recv1 := func(m *Message, reply SendFun) {
		if a1ReceivedIndex == 0 {
			require.Equal(t, a0InitMessage, string(m.Payload))
		} else {
			require.Equal(t, a0Reply, string(m.Payload))
		}
		a1ReceivedIndex++
		reply(SimpleMessage(a1Reply))
	}
	conns, agents, stopChat := NewChatSetup(recv0, recv1)
	defer stopChat()

	// kick-off of the conversation by A0
	// A1: expect the request from A0
	// A0: expect the reply from A1
	conns[1].ExpectMessages(1)
	conns[0].ExpectMessages(1)
	agents[0].Say(SimpleMessage(a0InitMessage))

	// wait for messages to arrive
	conns[1].WaitForMessages()
	conns[0].WaitForMessages()

	conns[1].ExpectMessages(1)
	agents[0].Say(SimpleMessage(a0Reply))
	conns[1].WaitForMessages()
}

func TestComplexTalk(t *testing.T) {
	// A0 -> A1 -> A2 chain of requests
	// A0 <- A1 <- A2 chain of replies with A2's birthday

	a0Request := "A1, please tell me the birthday of A2."
	a1Request := "A2, when is your birthday?"
	a2Reply := "01.01.2021"

	receivedBirthday := ""
	recv12 := func(m *Message, reply SendFun) {
		// save birthday reply from A2
		receivedBirthday = string(m.Payload)
		require.Equal(t, a2Reply, receivedBirthday)
	}
	recv2 := func(m *Message, reply SendFun) {
		require.Equal(t, a1Request, string(m.Payload))

		// reply to A1
		reply(SimpleMessage(a2Reply))
	}
	conns12, agents12, stopChat12 := NewChatSetup(recv12, recv2)
	defer stopChat12()

	recv0 := func(m *Message, reply SendFun) {
		// check received birthday value
		require.Equal(t, a2Reply, string(m.Payload))
	}
	recv10 := func(m *Message, reply SendFun) {
		require.Equal(t, a0Request, string(m.Payload))

		// forward request to A2
		// A1: expect the request from A0
		// A0: expect the final reply from A1
		conns12[1].ExpectMessages(1)
		conns12[0].ExpectMessages(1)
		agents12[0].Say(SimpleMessage(a1Request))

		// wait for messages to arrive
		conns12[1].WaitForMessages()
		conns12[0].WaitForMessages()

		// forward A2 reply to A0
		reply(SimpleMessage(receivedBirthday))
	}
	conns01, agents01, stopChat01 := NewChatSetup(recv0, recv10)
	defer stopChat01()

	// kick-off of the conversation by A0
	// A1: expect the request from A0
	// A0: expect the final reply from A1
	conns01[1].ExpectMessages(1)
	conns01[0].ExpectMessages(1)
	agents01[0].Say(SimpleMessage(a0Request))

	// wait for messages to arrive
	conns01[1].WaitForMessages()
	conns01[0].WaitForMessages()
}
