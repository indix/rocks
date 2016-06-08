package ops

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tecbot/gorocksdb"
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

	log.Printf("Trying to restore backup from %s to %s and WAL is going to %s\n", source, destination, walDestinationDir)
	return DoRestore(source, destination, walDestinationDir)
}

// DoRestore triggers a restore from the specified backup location
func DoRestore(source, destination, walDestinationDir string) error {
	opts := gorocksdb.NewDefaultOptions()
	db, err := gorocksdb.OpenBackupEngine(opts, source)
	if err != nil {
		return err
	}
	return db.RestoreDBFromLatestBackup(destination, walDestinationDir, gorocksdb.NewRestoreOptions())
}

func init() {
	Rocks.AddCommand(restore)

	restore.PersistentFlags().StringVar(&source, "src", "", "Restore from")
	restore.PersistentFlags().StringVar(&destination, "dest", "", "Restore to")
	restore.PersistentFlags().StringVar(&walDestinationDir, "wal", "", "Restore WAL to (generally same as --dest)")
}
