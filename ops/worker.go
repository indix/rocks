package ops

import (
	"log"
	"sync"
)

// WorkRequest struct contains source and destination for backup / restore
type WorkRequest struct {
	Source      string
	Destination string
}

// Worker for now
type Worker struct {
	running bool
	Queue   chan WorkRequest
	Errs    chan error
	Op      func(WorkRequest) error
	Marker  *sync.WaitGroup
}

// Start a worker
func (w *Worker) Start() {
	log.Printf("Starting Worker..\n")
	w.running = true
	go w.run()
}

func (w *Worker) run() {
	for w.running {
		select {
		case work := <-w.Queue:
			log.Printf("[Worker] Got work as %v\n", work)
			if err := w.Op(work); err != nil {
				w.Errs <- err
			}
			w.Marker.Done()
		}
	}
}

// Stop a worker
func (w *Worker) Stop() {
	w.running = false
	log.Printf("Stopped Worker..\n")
}
