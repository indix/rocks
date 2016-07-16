package consistency

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ind9/rocks/cmd/statistics"
	"github.com/spf13/cobra"
)

var consistencySource string
var consistencyRestore string
var counterFlag = 0
var recursive bool

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
		flagCheck, err = DoRecursiveConsistency(consistencySource, consistencyRestore)
	} else {
		err = DoConsistency(consistencySource, consistencyRestore)
	}

	if flagCheck == 0 {
		fmt.Printf("\nPASS: Source directory: %s and it's Restore: %s are consistant\n", consistencySource, consistencyRestore)
	}
	return err
}

// DoRecursiveConsistency checks for consistency recursively
func DoRecursiveConsistency(source, restore string) (int, error) {
	log.Printf("Initializing consistency check between %s data directory and %s as it's restore directory\n", source, restore)

	err := filepath.Walk(source, func(path string, info os.FileInfo, walkErr error) error {
		if info.Name() == Current {
			sourceDbLoc := filepath.Dir(path)
			sourceDbRelative, err := filepath.Rel(source, sourceDbLoc)
			restoreDbLoc := filepath.Join(restore, sourceDbRelative)

			if err = DoConsistency(sourceDbLoc, restoreDbLoc); err != nil {
				return err
			}
			return filepath.SkipDir
		}
		return walkErr
	})
	return counterFlag, err
}

// DoConsistency checks for consistency between rocks source store and its restore
func DoConsistency(source, restore string) error {
	log.Printf("Initializing consistency check between %s rocks store and %s as it's restore\n", source, restore)

	var rowCountSource, rowCountRestore int64
	log.Printf("Trying to collect the stores with non-matching number of keys\n")

	rowCountSource, err := statistics.DoStats(source)
	rowCountRestore, err = statistics.DoStats(restore)

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
}
