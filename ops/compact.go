package ops

import (
	"fmt"

	"github.com/spf13/cobra"
)

var compactionSource string
var compactionDestination string
var compactionThreads int

var backup = &cobra.Command{
	Use:   "compact",
	Short: "Does a compaction on rocksdb stores",
	Long:  "Does a compaction on rocksdb stores",
	Run:   AttachHandler(compactDatabase),
}

func compactDatabase(args []string) error {
	if compactionSource == "" {
		return fmt.Errorf("--src was not set")
	}
	if compactionDestination == "" {
		return fmt.Errorf("--dest was not set")
	}
	if recursive {
		return DoRecursiveCompaction(compactionSource, compactionDestination, compactionThreads)
	}
	return DoCompaction(compactionSource, compactionDestination)
}
