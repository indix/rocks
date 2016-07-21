package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/ind9/rocks/cmd/testutils"

	"github.com/spf13/cobra"
)

// Current is to identify a rocksdb store.
const Current = "CURRENT"

// LatestBackup is used to find the backup location
const LatestBackup = "LATEST_BACKUP"

var logDestination string // generally same as current working directory

// Rocks is the entry point command in the application
var Rocks = &cobra.Command{
	Use:   "rocks",
	Short: "RocksDB Ops CLI",
	Long: `Perform common ops related tasks on one or many RocksDB instances.

Find more details at https://github.com/ind9/rocks`,
}

func init() {
	Rocks.PersistentFlags().StringVar(&logDestination, "logDest", "", "Write logs to (generally same as current working directory)")
}

// CommandHandler is the wrapper interface that all commands to be implement as part of their "Run"
type CommandHandler func(args []string) error

func createLogs() error {
	var err error

	if logDestination == "" {
		logDestination, err = os.Getwd()
	}

	if !testutils.Exists(logDestination) {
		os.MkdirAll(logDestination, os.ModePerm)
	}

	var pathForLog = filepath.Join(logDestination, "rocks_"+time.Now().Format(time.RFC850)+".log")
	file, err := os.Create(pathForLog)
	defer file.Close()
	logFile, err := os.OpenFile(pathForLog, os.O_RDWR, 0644)
	if err != nil {
		log.Printf("error opening log file: %v", err)
		log.Printf("Logging on terminal . . . \n")
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(logFile)
	}
	defer logFile.Close()
	return err
}

// AttachHandler is a wrapper method for all commands that needs to be exposed
func AttachHandler(handler CommandHandler) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		start := time.Now()
		err := createLogs()
		err = handler(args)
		elapsed := time.Since(start).Seconds()
		fmt.Printf("This took  %f seconds\n", elapsed)
		if err != nil {
			log.Printf("[Error] %s", err.Error())
			os.Exit(1)
		}
	}
}
