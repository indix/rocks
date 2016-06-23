package ops

import "github.com/hashicorp/go-multierror"

// WorkerPool is a wrapper to manage a set of Workers efficiently
type WorkerPool struct {
	MaxWorkers int
	Op         func(WorkRequest) error
	workers    []Worker
	items      chan WorkRequest
	errs       chan error
	finalError error
}

// Initialize the workerpool
func (pool *WorkerPool) Initialize() {
	pool.items = make(chan WorkRequest)
	pool.errs = make(chan error)
	// Error handler
	go func(combined *error) {
		for err := range pool.errs {
			multierror.Append(*combined, err)
		}
	}(&pool.finalError)
}

// AddWork to a worker in the Pool
func (pool *WorkerPool) AddWork(work WorkRequest) {
	if len(pool.workers) < pool.MaxWorkers {
		worker := Worker{
			Queue: pool.items,
			Errs:  pool.errs,
			Op:    pool.Op,
		}
		worker.Start()
		pool.workers = append(pool.workers, worker)
	}
	pool.items <- work
}

// Join waits for all the tasks to complete - pool is not usable after this
func (pool *WorkerPool) Join() error {
	close(pool.items)
	close(pool.errs)
	return pool.finalError
}
