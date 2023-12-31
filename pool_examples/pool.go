package poolexamples

import (
	"sync"
	"sync/atomic"
)

type pool interface {
	execute(func())
	wait()
	activeCount() int32
}

type workerPool struct {
	activeWorkersCount atomic.Int32
	workers            chan struct{}
	wg                 *sync.WaitGroup
}

func newPool(limit int) pool {
	return &workerPool{
		workers: make(chan struct{}, limit),
		wg:      &sync.WaitGroup{},
	}
}

func (p *workerPool) execute(fn func()) {
	p.wg.Add(1)
	go func() {
		// block until there is a seat available in the pool
		p.workers <- struct{}{}
		p.activeWorkersCount.Add(1)

		defer p.wg.Done()
		fn()
		<-p.workers
		p.activeWorkersCount.Add(-1)
	}()
}

func (p *workerPool) wait() {
	p.wg.Wait()
}

func (p *workerPool) activeCount() int32 {
	return p.activeWorkersCount.Load()
}
