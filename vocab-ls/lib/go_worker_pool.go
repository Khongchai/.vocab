package lib

import (
	"context"
	"runtime"
	"sync"
)

// CPU-bound worker pool for cpu-intensive tasks.
type GoWorkerPool struct {
	ctx           context.Context
	wg            sync.WaitGroup
	channel       chan struct{}
	resourceMutex sync.Map
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
		pool.workWithMutex(resource, work)
	}()
}

func (pool *GoWorkerPool) WaitAll() {
	pool.wg.Wait()
}

func (pool *GoWorkerPool) workWithMutex(res string, work func()) {
	got, _ := pool.resourceMutex.LoadOrStore(res, &sync.Mutex{})
	mutex := got.(*sync.Mutex)

	mutex.Lock()
	defer mutex.Unlock()
	work()
}
