package ops

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tecbot/gorocksdb"
)

var consistancySource string
var consistancyRestore string
var flag bool

var consistency = &cobra.Command{
	Use:   "consistency",
	Short: "Checks for consistency between rocks store and it's corresponding restore",
	Long:  "Checks for the consistency between rocks store and it's corresponding restore",
	Run:   AttachHandler(checkConsistency),
}

func checkConsistency(args []string) error {
	if consistancySource == "" {
		return fmt.Errorf("--src was not set")
	}

	if consistancyRestore == "" {
		return fmt.Errorf("--dest was not set")
	}

	return DoConsistency(consistancySource, consistancyRestore)
}

// DoConsistency checks for consistency between rocks source store and its restore
func DoConsistency(source, restore string) error {
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
		fmt.Printf("Store : %s and corresponding Restore %s number of keys did not match\n", source, restore)
		fmt.Printf("Store Count : %v\n", rowCountSource)
		fmt.Printf("Restore Count : %v\n", rowCountRestore)
		flag = false
	}
	if err != nil {
		return err
	}
	return nil
}
