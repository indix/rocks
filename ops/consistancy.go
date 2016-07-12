package ops

import (
	"fmt"

	"github.com/spf13/cobra"
)

var consistancySource string
var consistancyRestore string

var consistency = &cobra.Command{
	Use:   "consistency",
	Short: "Checks for consistency between rocks store and it's corresponding restore",
	Long:  "Checks for the consistency between rocks store and it's corresponding restore",
	Run:   AttachHandler(checkConsistency),
}

func checkConsistency(args []string) error {
	if consistancySource == "" {
		return fmt.Errorf("--src was not set")
	}

	if consistancyRestore == "" {
		return fmt.Errorf("--dest was not set")
	}

	return DoConsistency(consistancySource, consistancyRestore)
}
