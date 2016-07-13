package restore

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hashicorp/go-multierror"
	"github.com/ind9/rocks/ops"
	"github.com/spf13/cobra"
	"github.com/tecbot/gorocksdb"
)

var restoreSource string
var restoreDestination string
var walDestinationDir string // generally the same as restoreDestination
var recursive bool
var keepLogFiles bool
var restoreThreads int

// LatestBackup is used to find the backup location
const LatestBackup = "LATEST_BACKUP"

var restore = &cobra.Command{
	Use:   "restore",
	Short: "Restore backed up rocksdb files",
	Long:  "Restore backed up rocksdb files",
	Run:   ops.AttachHandler(restoreDatabase),
}

func restoreDatabase(args []string) error {
	if restoreSource == "" {
		return fmt.Errorf("--src was not set")
	}
	if restoreDestination == "" {
		return fmt.Errorf("--dest was not set")
	}

	if walDestinationDir == "" {
		walDestinationDir = restoreDestination
	}

	if recursive {
		return DoRecursiveRestore(restoreSource, restoreDestination, walDestinationDir, restoreThreads, keepLogFiles)
	}
	return DoRestore(restoreSource, restoreDestination, walDestinationDir, keepLogFiles)
}

// DoRecursiveRestore recursively restores a rocksdb keeping the folder structure intact from backup to restore location
func DoRecursiveRestore(source, destination, walDestinationDir string, numThreads int, keepLogFiles bool) error {
	workerPool := ops.WorkerPool{
		MaxWorkers: restoreThreads,
		Op: func(request ops.WorkRequest) error {
			work := request.(ops.RestoreWork)
			return DoRestore(work.Source, work.Destination, work.WalDir, keepLogFiles)
		},
	}
	workerPool.Initialize()

	err := filepath.Walk(source, func(path string, info os.FileInfo, walkErr error) error {
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

			work := ops.RestoreWork{
				Source:      dbLoc,
				Destination: dbRestoreLoc,
				WalDir:      walRestoreLoc,
			}
			workerPool.AddWork(work)
			return filepath.SkipDir
		}

		return walkErr
	})

	var result error
	if errFromWorkers := workerPool.Join(); errFromWorkers != nil {
		result = multierror.Append(result, errFromWorkers)
	}

	if err != nil {
		result = multierror.Append(result, err)
	}

	return result
}

// DoRestore triggers a restore from the specified backup location
func DoRestore(source, destination, walDestinationDir string, keepLogFiles bool) error {
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
	db.Close()
	log.Printf("Restore complete from %s to %s and WAL went into %s\n", source, destination, walDestinationDir)
	return err
}

func init() {
	ops.Rocks.AddCommand(restore)

	restore.PersistentFlags().StringVar(&restoreSource, "src", "", "Restore from")
	restore.PersistentFlags().StringVar(&restoreDestination, "dest", "", "Restore to")
	restore.PersistentFlags().StringVar(&walDestinationDir, "wal", "", "Restore WAL to (generally same as --dest)")
	restore.PersistentFlags().BoolVar(&recursive, "recursive", false, "Trying restoring in recursive fashion from src to dest")
	restore.PersistentFlags().BoolVar(&keepLogFiles, "keep-log-files", false, "If true, restore won't overwrite the existing log files in wal_dir")
	restore.PersistentFlags().IntVar(&restoreThreads, "threads", 2*runtime.NumCPU(), "Number of threads while restoring")
}
