package commands

import (
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/internal/input"
	"github.com/spf13/cobra"
)

var (
	statCmd = &cobra.Command{
		Use:   "stat",
		Short: "Print info about events",
		Run:   statRun,
	}
)

func init() {
}

func statRun(cmd *cobra.Command, args []string) {
	fmt.Println("StatRun")
	input.DoLoginIfNeeded()
}
