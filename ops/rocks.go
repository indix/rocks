package ops

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/spf13/cobra"
)

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

// AttachHandler is a wrapper method for all commands that needs to be exposed
func AttachHandler(handler CommandHandler) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		if logDestination == "" {
			var err error
			logDestination, err = os.Getwd()
			log.Printf(logDestination)
			log.Printf("[Error] %v", err)
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

		start := time.Now()
		err = handler(args)
		elapsed := time.Since(start).Seconds()
		fmt.Printf("This took  %f seconds\n", elapsed)
		if err != nil {
			log.Printf("[Error] %s", err.Error())
			os.Exit(1)
		}
	}
}
