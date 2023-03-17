package work

import "sync"

// Worker must be implemented by types that want to use the work pool
type Worker interface {
	Task()
}

// Pool provides a pool of goroutines to execute Worker tasks
type Pool struct {
	work chan Worker    // single unbuffered channel
	wg   sync.WaitGroup // single wait group
}

// New creates a new work pool
func New(maxGoroutines int) *Pool {
	p := Pool{
		work: make(chan Worker),
	}

	p.wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		// Spawn the worker goroutines
		go func() {
			for w := range p.work { // blocks until something appears in channel
				w.Task()
			}
			p.wg.Done() // only called when channel closes
		}() // goroutine executed right away
	}
	return &p
}

// Run submits work to the pool
func (p *Pool) Run(w Worker) {
	p.work <- w
}

// Shutdown waits for all the goroutines to shutdown
func (p *Pool) Shutdown() {
	close(p.work)
	p.wg.Wait()
}
