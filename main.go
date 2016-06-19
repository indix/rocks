package main

import (
	"fmt"
	"os"

	"github.com/ind9/rocks/ops"
)

// Version of the app
var Version = "dev-build"

func main() {
	if err := ops.Rocks.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
