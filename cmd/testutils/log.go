package testutils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ind9/rocks/cmd"
)

// CreateLogs function generates logs in a log file
func CreateLogs() error {
	var err error

	if cmd.LogDestination == "" {
		cmd.LogDestination, err = os.Getwd()
	}

	if !Exists(cmd.LogDestination) {
		os.MkdirAll(cmd.LogDestination, os.ModePerm)
	}

	var pathForLog = filepath.Join(cmd.LogDestination, "rocks_"+strings.Replace(time.Stamp, " ", "", -1)+".log")

	os.Create(pathForLog)

	logFile, err := os.OpenFile(pathForLog, os.O_RDWR|os.O_WRONLY, 0644)

	if err != nil {
		log.Printf("error opening log file: %v", err)
		log.Printf("Logging on terminal . . . \n")
		log.SetOutput(os.Stdout) //setting the output as terminal in this case
	} else {
		log.SetOutput(logFile)
	}
	return err
}
