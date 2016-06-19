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
var record bool
var keepLogFiles bool

// LatestBackup is used to find the backup location
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
		return walkBackupDir(source, destination, walDestinationDir)
	}
	return DoRestore(source, destination, walDestinationDir)
}

func walkBackupDir(source, destination, walDestinationDir string) error {
	return filepath.Walk(source, func(path string, info os.FileInfo, walkErr error) error {
		if info.Name() == LatestBackup {
			dbLoc := filepath.Dir(path)
			dbRelative, err := filepath.Rel(source, dbLoc)
			if err != nil {
				log.Print(err)
				return err
			}

			dbRestoreLoc := filepath.Join(destination, dbRelative)
			walRestoreLoc := filepath.Join(walDestinationDir, dbRelative)

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
	log.Printf("Trying to restore backup from %s to %s and WAL is going to %s\n", source, destination, walDestinationDir)

	opts := gorocksdb.NewDefaultOptions()
	db, err := gorocksdb.OpenBackupEngine(opts, source)
	if err != nil {
		return err
	}

	restoreOpts := gorocksdb.NewRestoreOptions()
	if keepLogFiles {
		restoreOpts.SetKeepLogFiles(1)
	}
	err = db.RestoreDBFromLatestBackup(destination, walDestinationDir, restoreOpts)
	log.Printf("Restore complete from %s to %s and WAL went into %s\n", source, destination, walDestinationDir)
	return err
}

func init() {
	Rocks.AddCommand(restore)

	restore.PersistentFlags().StringVar(&source, "src", "", "Restore from")
	restore.PersistentFlags().StringVar(&destination, "dest", "", "Restore to")
	restore.PersistentFlags().StringVar(&walDestinationDir, "wal", "", "Restore WAL to (generally same as --dest)")
	restore.PersistentFlags().BoolVar(&recursive, "recursive", false, "Trying restoring in recursive fashion from src to dest")
	restore.PersistentFlags().BoolVar(&keepLogFiles, "keep-log-files", false, "If true, restore won't overwrite the existing log files in wal_dir")
}
