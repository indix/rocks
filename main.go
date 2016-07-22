package main

import (
	"fmt"
	"os"

	"github.com/ind9/rocks/cmd"
	_ "github.com/ind9/rocks/cmd/backup"
	_ "github.com/ind9/rocks/cmd/consistency"
	_ "github.com/ind9/rocks/cmd/restore"
	_ "github.com/ind9/rocks/cmd/statistics"
	"github.com/ind9/rocks/cmd/testutils"
	_ "github.com/ind9/rocks/cmd/trigger"
)

// Version of the app
var Version = "dev-build"

func main() {
	if err := testutils.CreateLogs(); err != nil {
		fmt.Println(err)
	}
	if err := cmd.Rocks.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
