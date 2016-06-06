package ops

import (
	"fmt"

	"github.com/spf13/cobra"
)

var source string
var destination string
var walDestinationDir string // generally the same as destination

var restore = &cobra.Command{
	Use:   "restore",
	Short: "Restore backed up rocksdb files",
	Long:  "Restore backed up rocksdb files",
	Run:   AttachHandler(restoreDatabase),
}

func restoreDatabase(args []string) error {
	if source == "" {
		return fmt.Errorf("--src was not set")
	}
	if destination == "" {
		return fmt.Errorf("--dest was not set")
	}

	if walDestinationDir == "" {
		walDestinationDir = destination
	}

	fmt.Printf("Source=%s, Dest=%s, WAL Dest=%s\n", source, destination, walDestinationDir)
	return nil
}

func init() {
	Rocks.AddCommand(restore)

	restore.PersistentFlags().StringVar(&source, "src", "", "Restore from")
	restore.PersistentFlags().StringVar(&destination, "dest", "", "Restore to")
	restore.PersistentFlags().StringVar(&walDestinationDir, "wal-dest", "", "Restore WAL to (generally same as --dest)")
}
