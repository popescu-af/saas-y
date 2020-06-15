package connection

import (
	"fmt"
	"sync"
	"time"

	"github.com/popescu-af/saas-y/pkg/log"
)

// ReadFn is the type of function listening for messages on a channel.
type ReadFn func() (*Message, error)

// WriteFn is the type of function sending messages on a channel.
type WriteFn func(*Message) error

// CloseFn is the type of function that closes a channel.
type CloseFn func()

// Channel is the structure containing a channel's read/write functions.
type Channel struct {
	Read  ReadFn
	Write WriteFn
	Close CloseFn
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

// FullDuplex is a full-duplex connection that takes a full-duplex endpoint
// and a channel. It handles messages arriving on the channel while periodically
// polling the endpoint for new messages to be dispached on the channel.
type FullDuplex struct {
	endpoint       FullDuplexEndpoint
	channel        *Channel
	pollingPeriod  time.Duration
	stopReading    chan bool
	stopProcessing chan bool
	wgStop         sync.WaitGroup
	lock           sync.Mutex
	isRunning      bool
}

// NewFullDuplex creates a new, inactive full-duplex connection.
// Call Run to activate run it.
func NewFullDuplex(endpoint FullDuplexEndpoint, channel *Channel, pollingPeriod time.Duration) *FullDuplex {
	return &FullDuplex{
		endpoint:       endpoint,
		channel:        channel,
		pollingPeriod:  pollingPeriod,
		stopReading:    make(chan bool),
		stopProcessing: make(chan bool),
	}
}

// Run is a blocking function that takes a full-duplex endpoint
// and a channel. It handles messages arriving on the channel while periodically
// polling the endpoint for new messages to be dispached on the channel.
func (f *FullDuplex) Run() error {
	f.lock.Lock()
	if f.isRunning {
		f.lock.Unlock()
		return fmt.Errorf("already running")
	}

	var wg sync.WaitGroup
	wg.Add(2)

	msgCh := make(chan *Message)
	errCh := make(chan error)

	// reader
	go func() {
		defer wg.Done()

		for {
			select {
			case <-f.stopReading:
				log.Info("external stop")
				return
			default:
				m, err := f.channel.Read()
				if err != nil {
					errCh <- err
					return
				}
				msgCh <- m
			}
		}
	}()

	var finalError error

	// processor
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(f.pollingPeriod)

		for {
			select {
			case <-f.stopProcessing:
				log.Info("external stop")
				return
			case msg := <-msgCh:
				if err := f.endpoint.ProcessMessage(msg, f.channel.Write); err != nil {
					log.ErrorCtx("failed to process message", log.Context{"error": err})
					return
				}
			case err := <-errCh:
				log.ErrorCtx("received error", log.Context{"error": err})
				finalError = err
				return
			case t := <-ticker.C:
				if err := f.endpoint.Poll(t, f.channel.Write); err != nil {
					log.ErrorCtx("failed to poll", log.Context{"error": err})
					return
				}
			}
		}
	}()

	f.isRunning = true
	f.wgStop.Add(1)
	f.lock.Unlock()
	wg.Wait()

	f.lock.Lock()
	f.isRunning = false
	f.lock.Unlock()

	f.wgStop.Done()
	return finalError
}

// Stop stops a full-duplex connection.
func (f *FullDuplex) Stop() error {
	f.lock.Lock()
	isRunning := f.isRunning
	f.lock.Unlock()

	if !isRunning {
		return fmt.Errorf("not running")
	}

	f.stopReading <- true
	f.stopProcessing <- true
	f.wgStop.Wait()

	f.channel.Close()
	return nil
}

// IsRunning returns the running status of a full-duplex connection.
func (f *FullDuplex) IsRunning() bool {
	f.lock.Lock()
	defer f.lock.Unlock()

	return f.isRunning
}
