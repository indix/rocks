package ops

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
)

var compactionSource string
var compactionDestination string
var compactionThreads int

var backup = &cobra.Command{
	Use:   "compact",
	Short: "Does a compaction on rocksdb stores",
	Long:  "Does a compaction on rocksdb stores",
	Run:   AttachHandler(compactDatabase),
}

func compactDatabase(args []string) error {
	if compactionSource == "" {
		return fmt.Errorf("--src was not set")
	}
	if compactionDestination == "" {
		return fmt.Errorf("--dest was not set")
	}
	if recursive {
		return DoRecursiveCompaction(compactionSource, compactionDestination, compactionThreads)
	}
	return DoCompaction(compactionSource, compactionDestination)
}

// DoRecursiveCompaction recursively compacts a rocksdb store keeping the folder structure intact as in source
func DoRecursiveCompaction(source, destination string, threads int) error {

	workerPool := WorkerPool{
		MaxWorkers: threads,
		Op: func(request WorkRequest) error {
			work := request.(CompactionWork)
			return DoCompaction(work.Source, work.Destination)
		},
	}
	workerPool.Initialize()

	err := filepath.Walk(source, func(path string, info os.FileInfo, walkErr error) error {
		if info.Name() == Current {
			dbLoc := filepath.Dir(path)

			dbRelative, err := filepath.Rel(source, dbLoc)
			if err != nil {
				log.Print(err)
				return err
			}

			dbCompactionLoc := filepath.Join(destination, dbRelative)
			if err = os.MkdirAll(dbCompactionLoc, os.ModePerm); err != nil {
				log.Print(err)
				return err
			}

			work := CompactionWork{
				Source:      dbLoc,
				Destination: dbCompactionLoc,
			}
			workerPool.AddWork(work)
			return filepath.SkipDir
		}
		return walkErr
	})

	var result error
	if errFromWorkers := workerPool.Join(); errFromWorkers != nil {
		result = multierror.Append(result, errFromWorkers)
	}

	if err != nil {
		result = multierror.Append(result, err)
	}

	return result
}
