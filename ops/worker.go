package ops

import "sync"

// WorkRequest struct contains source and destination for backup / restore
type WorkRequest struct {
	Source      string
	Destination string
}

// Worker for now
type Worker struct {
	Queue       chan WorkRequest
	Errs        chan error
	Op          func(WorkRequest) error
	ItemsMarker *sync.WaitGroup
	ErrsMarker  *sync.WaitGroup
}

// Start a worker
func (w *Worker) Start() {
	go w.run()
}

func (w *Worker) run() {
	for work := range w.Queue {
		if err := w.Op(work); err != nil {
			w.ErrsMarker.Add(1)
			w.Errs <- err
		}
	}
}
