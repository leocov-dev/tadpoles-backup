package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	backupCmd = &cobra.Command{
		Use:   "backup",
		Short: "Backup New Images.",
		Run:   backupRun,
	}
)

func init() {
}

func backupRun(cmd *cobra.Command, args []string) {
	fmt.Println("Backup Started ...")
	fmt.Println("*** NOT IMPLEMENTED ***")
}
