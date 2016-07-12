package ops

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tecbot/gorocksdb"
)

var consistancySource string
var consistancyRestore string
var consistancyThreads int
var flag bool

var consistency = &cobra.Command{
	Use:   "consistency",
	Short: "Checks for consistency between rocks store and it's corresponding restore",
	Long:  "Checks for the consistency between rocks store and it's corresponding restore",
	Run:   AttachHandler(checkConsistency),
}

func checkConsistency(args []string) (err error) {
	if consistancySource == "" {
		return fmt.Errorf("--src was not set")
	}

	if consistancyRestore == "" {
		return fmt.Errorf("--dest was not set")
	}
	var checkConsistant bool
	if recursive {
		checkConsistant, err = DoRecursiveConsistency(consistancySource, consistancyRestore, consistancyThreads)
	} else {
		checkConsistant, err = DoConsistency(consistancySource, consistancyRestore)
	}
	if checkConsistant {
		fmt.Printf("Store and Restore are consistant")
	}
	return err
}

// DoRecursiveConsistency checks for consistency recursively
func DoRecursiveConsistency(source, restore string, threads int) (bool, error) {
	/*
	   Implement this piece of code
	*/
	return false, nil
}

// DoConsistency checks for consistency between rocks source store and its restore
func DoConsistency(source, restore string) (bool, error) {
	log.Printf("Initializing consistency check between %s rocks store and %s as it's restore\n", source, restore)

	opts := gorocksdb.NewDefaultOptions()
	dbSource, err := gorocksdb.OpenDb(opts, source)
	defer dbSource.Close()
	dbRestore, err := gorocksdb.OpenDb(opts, restore)
	defer dbRestore.Close()

	consistencyOpts := gorocksdb.NewDefaultReadOptions()
	consistencyOpts.SetFillCache(false)
	sourceIterator := dbSource.NewIterator(consistencyOpts)
	restoreIterator := dbRestore.NewIterator(consistencyOpts)

	var rowCountSource, rowCountRestore int64
	log.Printf("Trying to collect the stores with non-matching number of keys\n")
	flag = true
	for sourceIterator.SeekToFirst(); sourceIterator.Valid(); sourceIterator.Next() {
		rowCountSource++
	}
	for restoreIterator.SeekToFirst(); restoreIterator.Valid(); restoreIterator.Next() {
		rowCountRestore++
	}
	if rowCountSource != rowCountRestore {
		log.Printf("Store : %s and corresponding Restore %s number of keys did not match\n", source, restore)
		log.Printf("Store Count : %v\n", rowCountSource)
		log.Printf("Restore Count : %v\n", rowCountRestore)
		flag = false
	}
	if err != nil {
		return flag, err
	}
	return flag, err
}
