package ops

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statsSource string
var statsThreads int

var stats = &cobra.Command{
	Use:   "statistics",
	Short: "Displays current statistics for a rocksdb store",
	Long:  "Displays current statistics for a rocksdb store",
	Run:   AttachHandler(generateStats),
}

func generateStats(args []string) error {
	if statsSource == "" {
		return fmt.Errorf("--src was not set")
	}
	if recursive {
		return DoRecursiveStats(statsSource, statsThreads)
	}
	return DoStats(statsSource)
}
