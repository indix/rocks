package ops

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tecbot/gorocksdb"
)

// Current is to identify a rocksdb store.
const Current = "CURRENT"

var backup = &cobra.Command{
	Use:   "backup",
	Short: "Backs up rocksdb stores",
	Long:  "Backs up rocksdb stores",
	Run:   AttachHandler(backupDatabase),
}

func backupDatabase(args []string) error {
	if source == "" {
		return fmt.Errorf("--src was not set")
	}
	if destination == "" {
		return fmt.Errorf("--dest was not set")
	}
	if recursive {
		walkSourceDir(source, destination)
		return nil
	}
	log.Printf("Trying to create backup from %s to %s\n", source, destination)
	err := DoBackup(source, destination)
	log.Printf("Backup from %s to %s completed\n", source, destination)
	return err
}

func walkSourceDir(source, destination string) error {
	return filepath.Walk(source, func(path string, info os.FileInfo, walkErr error) error {

		if info.Name() == Current {
			dbLoc := filepath.Dir(path)
			dbRelative, err := filepath.Rel(source, dbLoc)
			if err != nil {
				log.Print(err)
				return err
			}

			dbBackupLoc := filepath.Join(destination, dbRelative)
			log.Printf("Backup at %s, would be stored to %s\n", dbLoc, dbBackupLoc)
			log.Printf("Backup is created for %s rocks store", dbRelative)

			if err = os.MkdirAll(dbBackupLoc, os.ModePerm); err != nil {
				log.Print(err)
				return err
			}

			if err = DoBackup(dbLoc, dbBackupLoc); err != nil {
				log.Print(err)
				return err
			}
			return filepath.SkipDir
		}
		return walkErr
	})
}

// DoBackup triggers a backup from the source
func DoBackup(source, destination string) error {

	opts := gorocksdb.NewDefaultOptions()
	db, err := gorocksdb.OpenDb(opts, source)
	if err != nil {
		return err
	}

	backup, err := gorocksdb.OpenBackupEngine(opts, destination)
	if err != nil {
		return err
	}
	return backup.CreateNewBackup(db)
}

func init() {
	Rocks.AddCommand(backup)

	backup.PersistentFlags().StringVar(&source, "src", "", "Backup from")
	backup.PersistentFlags().StringVar(&destination, "dest", "", "Backup to")
	backup.PersistentFlags().BoolVar(&recursive, "recursive", false, "Trying to backup in recursive fashion from src to dest")
}
