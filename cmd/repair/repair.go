package repair

import (
	"fmt"
	"log"

	"github.com/ind9/rocks/cmd"
	"github.com/spf13/cobra"
	"github.com/tecbot/gorocksdb"
)

var recursive bool
var source string

// Work struct contains source and destination for backup
type Work struct {
	Source string
}

var repair = &cobra.Command{
	Use:   "repair",
	Short: "Does a DBRepair. It re-writes all manifest files based on the SST files.",
	Long:  "Does a DBRepair. It re-writes all manifest files based on the SST files.",
	Run:   cmd.AttachHandler(repairDatabase),
}

func repairDatabase(args []string) error {
	if source == "" {
		return fmt.Errorf("--src was not set")
	}

	return DoRepair(source)
}

// DoRepair triggers a compaction from the source
func DoRepair(source string) error {
	log.Printf("Trying to repair db at %s\n", source)

	opts := gorocksdb.NewDefaultOptions()
	err := gorocksdb.RepairDb(source, opts)

	if err != nil {
		return err
	}

	log.Printf("Repair for %s completed\n", source)
	return err
}

func init() {
	cmd.Rocks.AddCommand(repair)

	repair.PersistentFlags().StringVar(&source, "src", "", "Repair the source RocksDB")
}
