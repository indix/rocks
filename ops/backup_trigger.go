package ops

import "github.com/spf13/cobra"

var remote = &cobra.Command{
	Use:   "remote",
	Short: "Triggers a backup on a remote system",
	Long:  `Triggers a backup on a remote system. Currently it supports only HTTP based triggers.`,
	Run:   AttachHandler(backupDatabase),
}

var customHeaders []string // given in key=value format via command line
var url string
var method string

func init() {
	backup.AddCommand(remote)
}
