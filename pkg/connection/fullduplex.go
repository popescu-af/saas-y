package connection

import (
	"fmt"
	"sync"
)

// ReaderFn is the type for functions listening for messages.
type ReaderFn func() (*Message, error)

// WriterFn is the type for functions sending messages.
type WriterFn func(*Message) error

// Channel is the structure containing a channel's read/write functions.
type Channel struct {
	Read  ReaderFn
	Write WriterFn
}

// Message is a structure encompassing a message's type and payload.
type Message struct {
	Type    int
	Payload []byte
}

// FullDuplexEndpoint is the type for communication endpoints.
// Communication endpoints function in a full-duplex way, being possible to
// receive and send messages on the channel at the same time.
// Processing received messages has priority. If there is nothing to process,
// the endpoint can do housekeeping work or even send messages to the other side
// of the communication channel.
type FullDuplexEndpoint interface {
	// ProcessMessage should process the given message and react accordingly,
	// by changing state and/or writing something back or none of them.
	ProcessMessage(*Message, WriterFn) error
	// NextIteration should inspect the current state and decide whether
	// to take any action and/or write something on the connections' channel or not.
	NextIteration(WriterFn) error
}

// HandleTwoWayConnection is a blocking function that takes a full-duplex endpoint
// and a channel given by its read/write functions and handles the communication
// between the two.
func HandleTwoWayConnection(endpoint FullDuplexEndpoint, channel Channel) {
	var wg sync.WaitGroup
	wg.Add(2)

	msgCh := make(chan *Message)
	errCh := make(chan error)

	// reader
	go func() {
		defer wg.Done()

		for {
			m, err := channel.Read()
			if err != nil {
				errCh <- err
				return
			}
			msgCh <- m
		}
	}()

	// processor
	go func() {
		defer wg.Done()

		for {
			select {
			case msg := <-msgCh:
				if err := endpoint.ProcessMessage(msg, channel.Write); err != nil {
					fmt.Printf("failed to process message: %v", err)
					return
				}
			case err := <-errCh:
				fmt.Printf("received error: %v", err)
				return
			default:
				if err := endpoint.NextIteration(channel.Write); err != nil {
					fmt.Printf("failed to iterate: %v", err)
					return
				}
			}
		}
	}()

	wg.Wait()
}
