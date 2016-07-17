package ops

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

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

func generateStats(args []string) (err error) {
	if statsSource == "" {
		return fmt.Errorf("--src was not set")
	}
	var count int64
	if recursive {
		count, err = DoRecursiveStats(statsSource, statsThreads)
	} else {
		count, err = DoStats(statsSource)
	}
	fmt.Printf("Row Count on %s is %d\n", statsSource, count)
	return err
}

// DoRecursiveStats recursively generates statistics for a rocksdb store keeping the folder structure intact as in source
func DoRecursiveStats(source string, threads int) (int64, error) {
	var countsChan = make(chan int64)
	var wg sync.WaitGroup
	var totalRecordsCount int64
	go func(countsChan <-chan int64, totalRecordsCount *int64) {
		wg.Add(1)
		for dbCount := range countsChan {
			*totalRecordsCount += dbCount
		}
		wg.Done()
	}(countsChan, &totalRecordsCount)

	workerPool := WorkerPool{
		MaxWorkers: threads,
		Op: func(request WorkRequest) error {
			work := request.(StatsWork)
			count, err := DoStats(work.Source)
			work.Count <- count
			return err
		},
	}
	workerPool.Initialize()

	err := filepath.Walk(source, func(path string, info os.FileInfo, walkErr error) error {
		if info.Name() == Current {
			dbLoc := filepath.Dir(path)

			work := StatsWork{
				Source: dbLoc,
				Count:  countsChan,
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

	close(countsChan)
	wg.Wait()

	return totalRecordsCount, result
}

// DoStats generates statistics for the source
func DoStats(source string) (int64, error) {
	log.Printf("Trying to generate statistics for %s\n", source)

	opts := gorocksdb.NewDefaultOptions()
	db, err := gorocksdb.OpenDb(opts, source)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	statsOpts := gorocksdb.NewDefaultReadOptions()
	statsOpts.SetFillCache(false)
	iterator := db.NewIterator(statsOpts)

	var rowCount int64

	for iterator.SeekToFirst(); iterator.Valid(); iterator.Next() {
		rowCount++
	}

	if err = iterator.Err(); err != nil {
		iterator.Close()
		return rowCount, err
	}
	iterator.Close()

	return rowCount, nil
}

func init() {
	Rocks.AddCommand(stats)

	stats.PersistentFlags().StringVar(&statsSource, "src", "", "Statistics for")
	stats.PersistentFlags().BoolVar(&recursive, "recursive", false, "Trying to generate statistics in recursive fashion for src")
	stats.PersistentFlags().IntVar(&statsThreads, "threads", 2*runtime.NumCPU(), "Number of threads to generate statistics")
}
