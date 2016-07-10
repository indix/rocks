package ops

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"github.com/tecbot/gorocksdb"
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

// DoStats generates statistics for the source
func DoStats(source string) error {
	log.Printf("Trying to generate statistics for %s\n", source)

	opts := gorocksdb.NewDefaultOptions()
	db, err := gorocksdb.OpenDb(opts, source)
	if err != nil {
		return err
	}

	statsOpts := gorocksdb.NewDefaultReadOptions()
	statsOpts.SetFillCache(false)
	iterator := db.NewIterator(statsOpts)

	for iterator.SeekToFirst(); iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		value := iterator.Value()
		fmt.Printf("Key : %v  Value : %v\n", key.Data(), value.Data())
		key.Free()
		value.Free()
	}

	if err = iterator.Err(); err != nil {
		return err
	}

	iterator.Close()
	db.Close()
	log.Printf("Statistics generated from source %s\n", source)
	return nil
}

func init() {
	Rocks.AddCommand(stats)

	stats.PersistentFlags().StringVar(&statsSource, "src", "", "Statistics for")
	stats.PersistentFlags().BoolVar(&recursive, "recursive", false, "Trying to generate statistics in recursive fashion for src")
	stats.PersistentFlags().IntVar(&statsThreads, "threads", 2*runtime.NumCPU(), "Number of threads to generate statistics")
}
