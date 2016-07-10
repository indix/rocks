package ops

import (
	"sync"
	"sync/atomic"

	"github.com/hashicorp/go-multierror"
)

// WorkerPool is a wrapper to manage a set of Workers efficiently
type WorkerPool struct {
	MaxWorkers  int
	Op          func(WorkRequest) error
	workers     []Worker
	items       chan WorkRequest
	itemsMarker sync.WaitGroup
	errs        chan error
	count       chan int64
	finalError  error
	finalCount  int64
}

// Initialize the workerpool
func (pool *WorkerPool) Initialize() {
	pool.items = make(chan WorkRequest)
	pool.errs = make(chan error)
	pool.count = make(chan int64)
	// Error handler
	go func(combined *error) {
		var result = *combined
		for err := range pool.errs {
			result = multierror.Append(result, err)
		}
		combined = &result
	}(&pool.finalError)
	// Counter handler
	go func(combined *int64) {
		var result = *combined
		for countForShard := range pool.count {
			result = atomic.AddInt64(&pool.finalCount, countForShard)
		}
		combined = &result
	}(&pool.finalCount)
}

// AddWork to a worker in the Pool
func (pool *WorkerPool) AddWork(work WorkRequest) {
	if len(pool.workers) < pool.MaxWorkers {
		worker := Worker{
			Queue:  pool.items,
			Errs:   pool.errs,
			Op:     pool.Op,
			Count:  pool.count,
			Marker: &pool.itemsMarker,
		}
		worker.Start()
		pool.workers = append(pool.workers, worker)
	}
	pool.itemsMarker.Add(1)
	pool.items <- work
}

// Join waits for all the tasks to complete - pool is not usable after this
func (pool *WorkerPool) Join() error {
	close(pool.items)
	pool.itemsMarker.Wait()

	close(pool.errs)
	return pool.finalError
}

// JoinCount returns total count of all the keys in all the shards - pool is not usable after this
func (pool *WorkerPool) JoinCount() int64 {
	close(pool.count)
	return atomic.LoadInt64(&pool.finalCount)
}
