package compact

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ashwanthkumar/golang-utils/worker"
	"github.com/hashicorp/go-multierror"
	"github.com/ind9/rocks/cmd"
	"github.com/spf13/cobra"
	"github.com/tecbot/gorocksdb"
)

var compactionSource string
var compactionThreads int
var keys gorocksdb.Range
var recursive bool

// Work struct contains source and destination for backup
type Work struct {
	Source string
}

var compact = &cobra.Command{
	Use:   "compact",
	Short: "Does a compaction on rocksdb stores",
	Long:  "Does a compaction on rocksdb stores",
	Run:   cmd.AttachHandler(compactDatabase),
}

func compactDatabase(args []string) error {
	if compactionSource == "" {
		return fmt.Errorf("--src was not set")
	}

	if recursive {
		return DoRecursiveCompaction(compactionSource, compactionThreads)
	}
	return DoCompaction(compactionSource)
}

// DoRecursiveCompaction recursively compacts a rocksdb store keeping the folder structure intact as in source
func DoRecursiveCompaction(source string, threads int) error {

	workerPool := worker.Pool{
		MaxWorkers: threads,
		Op: func(request worker.Request) error {
			work := request.(Work)
			return DoCompaction(work.Source)
		},
	}
	workerPool.Initialize()

	err := filepath.Walk(source, func(path string, info os.FileInfo, walkErr error) error {
		if info.Name() == cmd.Current {
			dbLoc := filepath.Dir(path)

			work := Work{
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
	cmd.Rocks.AddCommand(compact)

	compact.PersistentFlags().StringVar(&compactionSource, "src", "", "Compact for the source rocksdb store")
	compact.PersistentFlags().BoolVar(&recursive, "recursive", false, "Trying to compact in recursive fashion for src")
	compact.PersistentFlags().IntVar(&compactionThreads, "threads", 2*runtime.NumCPU(), "Number of threads to do backup")
}
