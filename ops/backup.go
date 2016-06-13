package ops

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tecbot/gorocksdb"
)

// MpidIndexDb is one of the rocksdb stores
const MpidIndexDb = "mpid_index_db"

// StoreDb is one of the rocksdb stores
const StoreDb = "store_db"

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
	return DoBackup(source, destination)
}

func walkSourceDir(source, destination string) {
	filepath.Walk(source, func(path string, info os.FileInfo, walkErr error) error {
		if info.Name() == MpidIndexDb || info.Name() == StoreDb {
			dbLoc := filepath.Dir(path)
			dbRelative, err := filepath.Rel(source, dbLoc)
			if err != nil {
				log.Print(err)
				return err
			}

			dbBackupLoc := filepath.Join(destination, dbRelative)
			log.Printf("Backup at %s, would be stored to %s\n", dbLoc, dbBackupLoc)

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
	files, err := ioutil.ReadDir(source)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		name := file.Name()

		if (name != MpidIndexDb) && (name != StoreDb) {
			continue
		}

		opts := gorocksdb.NewDefaultOptions()
		rocksdb, err := gorocksdb.OpenDb(opts, filepath.Join(source, name))
		db, err := gorocksdb.OpenBackupEngine(opts, destination)

		if err != nil {
			return err
		}
		db.CreateNewBackup(rocksdb)
	}
	return nil
}

func init() {
	Rocks.AddCommand(backup)

	backup.PersistentFlags().StringVar(&source, "src", "", "Backup from")
	backup.PersistentFlags().StringVar(&destination, "dest", "", "Backup to")
	backup.PersistentFlags().BoolVar(&recursive, "recursive", false, "Trying to backup in recursive fashion from src to dest")
}
