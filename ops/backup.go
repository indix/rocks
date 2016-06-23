package ops

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"github.com/tecbot/gorocksdb"
)

// Current is to identify a rocksdb store.
const Current = "CURRENT"

var backup = &cobra.Command{
	Use:   "backup",
	Short: "Backs up rocksdb stores",
	Long:  "Backs up rocksdb stores",
	Run:   AttachHandler(backupDatabase),
}

func backupDatabase(args []string) error {
	if source == "" {
		return fmt.Errorf("--src was not set")
	}
	if destination == "" {
		return fmt.Errorf("--dest was not set")
	}
	if recursive {
		return walkSourceDir(source, destination)
	}
	return DoBackup(source, destination)
}

func walkSourceDir(source, destination string) error {
	workQueue := make(chan WorkRequest)
	errsQueue := make(chan error)
	var marker sync.WaitGroup
	var workers []Worker
	for workerCount := 0; workerCount < 5; workerCount++ {
		worker := Worker{
			Queue:  workQueue,
			Errs:   errsQueue,
			Marker: &marker,
			Op: func(request WorkRequest) error {
				log.Printf("Got work %v\n", request)
				return DoBackup(request.Source, request.Destination)
			},
		}

		worker.Start()
		workers = append(workers, worker)
	}

	err := filepath.Walk(source, func(path string, info os.FileInfo, walkErr error) error {
		if info.Name() == Current {
			dbLoc := filepath.Dir(path)

			dbRelative, err := filepath.Rel(source, dbLoc)
			if err != nil {
				log.Print(err)
				return err
			}

			dbBackupLoc := filepath.Join(destination, dbRelative)
			if err = os.MkdirAll(dbBackupLoc, os.ModePerm); err != nil {
				log.Print(err)
				return err
			}

			workRequest := WorkRequest{
				Source:      dbLoc,
				Destination: dbBackupLoc,
			}
			workQueue <- workRequest
			marker.Add(1)
			return filepath.SkipDir
		}
		return walkErr
	})
	if err != nil {
		return err
	}
	marker.Wait() // wait for all the items to get processed
	close(workQueue)
	close(errsQueue)
	for _, worker := range workers {
		worker.Stop()
	}

	var result error
	for item := range errsQueue {
		multierror.Append(result, item)
	}

	return err
}

// DoBackup triggers a backup from the source
func DoBackup(source, destination string) error {
	log.Printf("Trying to create backup from %s to %s\n", source, destination)

	opts := gorocksdb.NewDefaultOptions()
	db, err := gorocksdb.OpenDb(opts, source)
	if err != nil {
		return err
	}

	backup, err := gorocksdb.OpenBackupEngine(opts, destination)
	if err != nil {
		return err
	}
	err = backup.CreateNewBackup(db)
	log.Printf("Backup from %s to %s completed\n", source, destination)
	return err
}

func init() {
	Rocks.AddCommand(backup)

	backup.PersistentFlags().StringVar(&source, "src", "", "Backup from")
	backup.PersistentFlags().StringVar(&destination, "dest", "", "Backup to")
	backup.PersistentFlags().BoolVar(&recursive, "recursive", false, "Trying to backup in recursive fashion from src to dest")
}
