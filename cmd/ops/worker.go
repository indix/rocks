package ops

import "sync"

// WorkRequest Base type of all work objects
type WorkRequest interface{}

// RestoreWork struct contains source, destination and WAL for restore
type RestoreWork struct {
	Source      string
	Destination string
	WalDir      string
}

// BackupWork struct contains source and destination for backup
type BackupWork struct {
	Source      string
	Destination string
}

// StatsWork struct contains source for generating statistics
type StatsWork struct {
	Source string
	Count  chan<- int64
}

// ConsistencyWork struct contains source and restore locations for comparison
type ConsistencyWork struct {
	Source  string
	Restore string
}

// Worker for now
type Worker struct {
	Queue  chan WorkRequest
	Errs   chan error
	Op     func(WorkRequest) error
	Marker *sync.WaitGroup
}

// Start a worker
func (w *Worker) Start() {
	go w.run()
}

func (w *Worker) run() {
	for work := range w.Queue {
		if err := w.Op(work); err != nil {
			w.Errs <- err
		}
		w.Marker.Done()
	}
}
