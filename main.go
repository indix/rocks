package main

import (
	"fmt"
	"os"

	_ "github.com/ind9/rocks/cmd/backup"
	_ "github.com/ind9/rocks/cmd/consistency"
	_ "github.com/ind9/rocks/cmd/restore"
	_ "github.com/ind9/rocks/cmd/statistics"
	_ "github.com/ind9/rocks/cmd/trigger"
)

// Version of the app
var Version = "dev-build"

func main() {
	if err := Rocks.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
