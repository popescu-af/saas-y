package worker

import (
	"sync"
)

// Job represents a unit of work.
type Job func()

// Pool is the abstraction for a worker pool.
type Pool struct {
	wg   sync.WaitGroup
	jobs chan Job
	exit []chan struct{}
}

func (p *Pool) worker(index int) {
	defer p.wg.Done()
	for {
		select {
		case j := <-p.jobs:
			j()
		case <-p.exit[index]:
			return
		}
	}
}

// Post is used for posting new jobs that have to be executed.
func (p *Pool) Post(j Job) {
	p.jobs <- j
}

// Stop stops the worker threads. The pool is unusable after calling Stop().
func (p *Pool) Stop() {
	close(p.jobs)
	for _, e := range p.exit {
		e <- struct{}{}
	}
	p.wg.Wait()
}

// NewPool creates a new worker pool with 'n' workers.
func NewPool(n int) *Pool {
	p := &Pool{
		jobs: make(chan Job),
		exit: make([]chan struct{}, n),
	}
	p.wg.Add(n)
	for i := 0; i < n; i++ {
		p.exit[i] = make(chan struct{})
		go p.worker(i)
	}
	return p
}
