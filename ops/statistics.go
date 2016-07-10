package ops

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
)

var statsSource string
var statsThreads int

var stats = &cobra.Command{
	Use:   "statistics",
	Short: "Displays current statistics for a rocksdb store",
	Long:  "Displays current statistics for a rocksdb store",
	Run:   AttachHandler(generateStats),
}

func generateStats(args []string) error {
	if statsSource == "" {
		return fmt.Errorf("--src was not set")
	}
	if recursive {
		return DoRecursiveStats(statsSource, statsThreads)
	}
	return DoStats(statsSource)
}

// DoRecursiveStats recursively generates statistics for a rocksdb store keeping the folder structure intact as in source
func DoRecursiveStats(source string, threads int) error {

	workerPool := WorkerPool{
		MaxWorkers: threads,
		Op: func(request WorkRequest) error {
			work := request.(StatsWork)
			return DoStats(work.Source)
		},
	}
	workerPool.Initialize()

	err := filepath.Walk(source, func(path string, info os.FileInfo, walkErr error) error {
		if info.Name() == Current {
			dbLoc := filepath.Dir(path)

			work := StatsWork{
				Source: dbLoc,
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
