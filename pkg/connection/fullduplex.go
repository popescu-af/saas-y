package connection

import (
	"sync"
	"time"

	"github.com/popescu-af/saas-y/pkg/log"
)

// ReadFn is the type of function listening for messages on a channel.
type ReadFn func() (*Message, error)

// WriteFn is the type of function sending messages on a channel.
type WriteFn func(*Message) error

// Channel is the structure containing a channel's read/write functions.
type Channel struct {
	Read  ReadFn
	Write WriteFn
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
	ProcessMessage(*Message, WriteFn) error
	// Poll should inspect the current state and decide whether to take any action
	// and/or write something on the channel or not.
	Poll(time.Time, WriteFn) error
}

// TalkFullDuplex is a blocking function that takes a full-duplex endpoint
// and a channel. It handles messages arriving on the channel while periodically
// polling the endpoint for new messages to be dispached on the channel.
func TalkFullDuplex(endpoint FullDuplexEndpoint, channel Channel, pollingPeriod time.Duration) error {
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

	var finalError error

	// processor
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(pollingPeriod)

		for {
			select {
			case msg := <-msgCh:
				if err := endpoint.ProcessMessage(msg, channel.Write); err != nil {
					log.ErrorCtx("failed to process message", log.Context{"error": err})
					return
				}
			case err := <-errCh:
				log.ErrorCtx("received error", log.Context{"error": err})
				finalError = err
				return
			case t := <-ticker.C:
				if err := endpoint.Poll(t, channel.Write); err != nil {
					log.ErrorCtx("failed to poll", log.Context{"error": err})
					return
				}
			}
		}
	}()

	wg.Wait()
	return finalError
}
