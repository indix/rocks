package consistency

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	log "github.com/Sirupsen/logrus"

	"github.com/ashwanthkumar/golang-utils/worker"
	"github.com/hashicorp/go-multierror"
	"github.com/ind9/rocks/cmd"
	"github.com/ind9/rocks/cmd/statistics"
	"github.com/spf13/cobra"
)

var consistencySource string
var consistencyRestore string
var consistencyThreads int
var recursive bool

// Work struct contains Source and Restore locations for checking consistency
type Work struct {
	Source  string
	Restore string
}

var consistency = &cobra.Command{
	Use:   "consistency",
	Short: "Checks for consistency between rocks store and it's corresponding restore",
	Long:  "Checks for the consistency between rocks store and it's corresponding restore",
	Run:   cmd.AttachHandler(checkConsistency),
}

func checkConsistency(args []string) (err error) {
	if consistencySource == "" {
		return fmt.Errorf("--src was not set")
	}

	if consistencyRestore == "" {
		return fmt.Errorf("--dest was not set")
	}

	if recursive {
		return DoRecursiveConsistency(consistencySource, consistencyRestore, consistencyThreads)
	}
	return DoConsistency(consistencySource, consistencyRestore)
}

// DoRecursiveConsistency checks for consistency recursively
func DoRecursiveConsistency(source, restore string, threads int) error {
	log.Printf("Initializing consistency check between %s data directory and %s as it's restore directory\n", source, restore)

	workerPool := worker.Pool{
		MaxWorkers: threads,
		Op: func(request worker.Request) error {
			work := request.(Work)
			return DoConsistency(work.Source, work.Restore)
		},
	}
	workerPool.Initialize()

	err := filepath.Walk(source, func(path string, info os.FileInfo, walkErr error) error {
		if info.Name() == cmd.Current {
			sourceDbLoc := filepath.Dir(path)
			sourceDbRelative, err := filepath.Rel(source, sourceDbLoc)
			if err != nil {
				return err
			}
			restoreDbLoc := filepath.Join(restore, sourceDbRelative)

			work := Work{
				Source:  sourceDbLoc,
				Restore: restoreDbLoc,
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

// DoConsistency checks for consistency between rocks source store and its restore
func DoConsistency(source, restore string) error {
	var rowCountSource, rowCountRestore int64

	rowCountSource, err := statistics.DoStats(source)
	rowCountRestore, err = statistics.DoStats(restore)

	if rowCountSource != rowCountRestore {
		log.Printf("Source : %s Restore : %s", source, restore)
		log.Printf("Store Count : %v\n", rowCountSource)
		log.Printf("Restore Count : %v\n", rowCountRestore)
		log.Printf("STATUS : FAIL")
	} else {
		log.Printf("Source : %s Restore : %s", source, restore)
		log.Printf("Store Count : %v\n", rowCountSource)
		log.Printf("Restore Count : %v\n", rowCountRestore)
		log.Printf("STATUS : PASS")
	}
	if err != nil {
		return err
	}

	return nil
}

func init() {
	cmd.Rocks.AddCommand(consistency)

	consistency.PersistentFlags().StringVar(&consistencySource, "src", "", "Rocks store location")
	consistency.PersistentFlags().StringVar(&consistencyRestore, "dest", "", "Restore location for Rocks store")
	consistency.PersistentFlags().BoolVar(&recursive, "recursive", false, "Trying to check consistency between rocks store and and it's restore")
	consistency.PersistentFlags().IntVar(&consistencyThreads, "threads", 2*runtime.NumCPU(), "Number of threads to do backup")
}
