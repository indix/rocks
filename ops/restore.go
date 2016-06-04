package ops

import "github.com/spf13/cobra"

var restore = &cobra.Command{
	Use:   "restore",
	Short: "Restore backed up rocksdb files",
	Long:  "Restore backed up rocksdb files",
}

func init() {
	Rocks.AddCommand(restore)
}
