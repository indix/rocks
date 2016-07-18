package ops

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"github.com/tecbot/gorocksdb"
)

var compactionSource string
var compactionThreads int
var keys gorocksdb.Range
var logDestination string // generally same as current working directory

var compact = &cobra.Command{
	Use:   "compact",
	Short: "Does a compaction on rocksdb stores",
	Long:  "Does a compaction on rocksdb stores",
	Run:   AttachHandler(compactDatabase),
}

func compactDatabase(args []string) error {
	if compactionSource == "" {
		return fmt.Errorf("--src was not set")
	}

	if recursive {
		return DoRecursiveCompaction(compactionSource, compactionThreads)
	}

	if logDestination == "" {
		var err error
		logDestination, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	return DoCompaction(compactionSource)
}

// DoRecursiveCompaction recursively compacts a rocksdb store keeping the folder structure intact as in source
func DoRecursiveCompaction(source string, threads int) error {

	workerPool := WorkerPool{
		MaxWorkers: threads,
		Op: func(request WorkRequest) error {
			work := request.(CompactionWork)
			return DoCompaction(work.Source)
		},
	}
	workerPool.Initialize()

	err := filepath.Walk(source, func(path string, info os.FileInfo, walkErr error) error {
		if info.Name() == Current {
			dbLoc := filepath.Dir(path)

			work := CompactionWork{
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

// DoCompaction triggers a compaction from the source
func DoCompaction(source string) error {

	logFile, err := os.OpenFile(filepath.Join(logDestination, time.Now().Format(time.RFC850), "compact.log"), os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Printf("error opening log file: %v", err)
		log.Printf("Logging on terminal . . . \n")
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(logFile)
	}
	defer logFile.Close()
	log.Printf("Trying to compact data for %s\n", source)

	opts := gorocksdb.NewDefaultOptions()
	compactOpts := gorocksdb.NewDefaultReadOptions()
	db, err := gorocksdb.OpenDb(opts, source)
	defer db.Close()
	keys.Start = db.NewIterator(compactOpts).Key().Data()
	db.CompactRange(keys)

	if err != nil {
		return err
	}

	log.Printf("Compaction for %s completed\n", source)
	return err
}

func init() {
	Rocks.AddCommand(compact)

	compact.PersistentFlags().StringVar(&compactionSource, "src", "", "Compact for the source rocksdb store")
	compact.PersistentFlags().StringVar(&logDestination, "logDest", "", "Generate logs to (generally same as current working directory)")
	compact.PersistentFlags().BoolVar(&recursive, "recursive", false, "Trying to compact in recursive fashion for src")
	compact.PersistentFlags().IntVar(&compactionThreads, "threads", 2*runtime.NumCPU(), "Number of threads to do backup")
}
