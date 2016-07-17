package ops

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
)

var consistencySource string
var consistencyRestore string
var consistencyThreads int
var counterFlag = 0

var consistency = &cobra.Command{
	Use:   "consistency",
	Short: "Checks for consistency between rocks store and it's corresponding restore",
	Long:  "Checks for the consistency between rocks store and it's corresponding restore",
	Run:   AttachHandler(checkConsistency),
}

func checkConsistency(args []string) (err error) {
	if consistencySource == "" {
		return fmt.Errorf("--src was not set")
	}

	if consistencyRestore == "" {
		return fmt.Errorf("--dest was not set")
	}

	var flagCheck int
	if recursive {
		flagCheck, err = DoRecursiveConsistency(consistencySource, consistencyRestore, consistencyThreads)
	} else {
		err = DoConsistency(consistencySource, consistencyRestore)
	}

	if flagCheck == 0 {
		fmt.Printf("\nPASS: Source directory: %s and it's Restore: %s are consistant\n", consistencySource, consistencyRestore)
	}
	return err
}

// DoRecursiveConsistency checks for consistency recursively
func DoRecursiveConsistency(source, restore string, threads int) (int, error) {
	log.Printf("Initializing consistency check between %s data directory and %s as it's restore directory\n", source, restore)

	workerPool := WorkerPool{
		MaxWorkers: threads,
		Op: func(request WorkRequest) error {
			work := request.(ConsistencyWork)
			return DoConsistency(work.Source, work.Restore)
		},
	}
	workerPool.Initialize()

	err := filepath.Walk(source, func(path string, info os.FileInfo, walkErr error) error {
		if info.Name() == Current {
			sourceDbLoc := filepath.Dir(path)
			sourceDbRelative, err := filepath.Rel(source, sourceDbLoc)
			if err != nil {
				return err
			}
			restoreDbLoc := filepath.Join(restore, sourceDbRelative)

			work := ConsistencyWork{
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

	return counterFlag, result
}

// DoConsistency checks for consistency between rocks source store and its restore
func DoConsistency(source, restore string) error {
	log.Printf("Initializing consistency check between %s rocks store and %s as it's restore\n", source, restore)

	var rowCountSource, rowCountRestore int64
	log.Printf("Trying to collect the stores with non-matching number of keys\n")

	rowCountSource, err := DoStats(source)
	rowCountRestore, err = DoStats(restore)

	if rowCountSource != rowCountRestore {
		log.Printf("Store : %s and corresponding Restore %s number of keys did not match\n", source, restore)
		log.Printf("Store Count : %v\n", rowCountSource)
		log.Printf("Restore Count : %v\n", rowCountRestore)
		counterFlag++
	}
	if err != nil {
		return err
	}
	log.Printf("Store: %s is consistent with restore: %s\n", source, restore)
	return nil
}

func init() {
	Rocks.AddCommand(consistency)

	consistency.PersistentFlags().StringVar(&consistencySource, "src", "", "Rocks store location")
	consistency.PersistentFlags().StringVar(&consistencyRestore, "dest", "", "Restore location for Rocks store")
	consistency.PersistentFlags().BoolVar(&recursive, "recursive", false, "Trying to check consistency between rocks store and and it's restore")
	consistency.PersistentFlags().IntVar(&consistencyThreads, "threads", 2*runtime.NumCPU(), "Number of threads to do backup")
}
