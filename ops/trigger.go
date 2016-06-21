package ops

import "github.com/spf13/cobra"

var trigger = &cobra.Command{
	Use:   "trigger",
	Short: "Triggers a backup on a remote system",
	Long:  `Triggers a backup on a remote system`,
}

func init() {
	Rocks.AddCommand(trigger)
}
