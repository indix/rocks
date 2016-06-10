package ops

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tecbot/gorocksdb"
)

var source string
var destination string
var walDestinationDir string // generally the same as destination
var recursive bool

// LatestBackup is the file at which terminate recursive lookups
const LatestBackup = "LATEST_BACKUP"

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

	if recursive {
		walkDir(source, destination, walDestinationDir)
		return nil
	}
	log.Printf("Trying to restore backup from %s to %s and WAL is going to %s\n", source, destination, walDestinationDir)
	return DoRestore(source, destination, walDestinationDir)
}

func walkDir(source, destination, walDestinationDir string) {
	filepath.Walk(source, func(path string, info os.FileInfo, walkErr error) error {
		if info.Name() == LatestBackup {
			dbLoc := filepath.Dir(path)
			dbRelative, err := filepath.Rel(source, dbLoc)
			if err != nil {
				log.Print(err)
				return err
			}

			dbRestoreLoc := filepath.Join(destination, dbRelative)
			walRestoreLoc := filepath.Join(walDestinationDir, dbRelative)
			log.Printf("Backup at %s, would be restored to %s with WAL to %s\n", dbLoc, dbRestoreLoc, walRestoreLoc)

			if err = os.MkdirAll(dbRestoreLoc, os.ModePerm); err != nil {
				log.Print(err)
				return err
			}

			if err = os.MkdirAll(walRestoreLoc, os.ModePerm); err != nil {
				log.Print(err)
				return err
			}

			if err = DoRestore(dbLoc, dbRestoreLoc, walRestoreLoc); err != nil {
				log.Print(err)
				return err
			}

			return filepath.SkipDir
		}

		return walkErr
	})
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
	restore.PersistentFlags().BoolVar(&recursive, "recursive", false, "Trying restoring in recursive fashion from src to dest")
}
