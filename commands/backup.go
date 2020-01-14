package commands

import (
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/internal/input"
	"github.com/spf13/cobra"
)

var (
	backupCmd = &cobra.Command{
		Use:   "backup",
		Short: "Backup new images.",
		Run:   backupRun,
	}
)

func init() {
}

func backupRun(cmd *cobra.Command, args []string) {
	fmt.Println("BackupRun")
	input.DoLoginIfNeeded()
}
