package lib

import (
	"context"
	"runtime"
	"sync"
)

// CPU-bound worker pool for cpu-intensive tasks.
type GoWorkerPool struct {
	ctx     context.Context
	wg      sync.WaitGroup
	channel chan struct{}
}

func NewGoWorkerPool(ctx context.Context) *GoWorkerPool {
	runtime.GOMAXPROCS(runtime.NumCPU())
	return &GoWorkerPool{
		ctx:     ctx,
		channel: make(chan struct{}, runtime.NumCPU()),
		wg:      sync.WaitGroup{},
	}
}

// Spawn a goroutine to perform work.
//
// Blocks until there are remaining workers to schedule work onto.
func (pool *GoWorkerPool) Run(resource string, work func()) {
	pool.channel <- struct{}{}
	pool.wg.Add(1)
	go func() {
		defer func() {
			pool.wg.Done()
			<-pool.channel
		}()
		work()
	}()
}

func (pool *GoWorkerPool) WaitAll() {
	pool.wg.Wait()
}
