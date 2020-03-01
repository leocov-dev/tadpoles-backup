package commands

import (
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/internal/tadpoles_api"
	"github.com/leocov-dev/tadpoles-backup/internal/user_input"
	"github.com/leocov-dev/tadpoles-backup/internal/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	backupCmd = &cobra.Command{
		Use:   "backup [target-directory]",
		Short: "Backup New Images.",
		Run:   backupRun,
		Args:  backupArgs(),
		PreRun: func(cmd *cobra.Command, args []string) {
			user_input.DoLoginIfNeeded()
		},
	}
)

func init() {
	rootCmd.AddCommand(backupCmd)
}

func backupArgs() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("[target-directory] argument missing")
		}
		return nil
	}
}

func backupRun(cmd *cobra.Command, args []string) {
	fmt.Println("Backup Started...")

	backupTarget := filepath.Clean(args[0])
	err := os.MkdirAll(backupTarget, os.ModePerm)
	if err != nil {
		utils.CmdFailed(cmd, err)
	}
	log.Debug("Backup to: ", backupTarget)

	info, err := tadpoles_api.GetAccountInfo()
	if err != nil {
		utils.CmdFailed(cmd, err)
	}

	fmt.Print("Checking Events...")
	log.Debug("") // newline for debug mode
	attachments, err := tadpoles_api.GetEventAttachmentData(info.FirstEvent, info.LastEvent)
	if err != nil {
		utils.CmdFailed(cmd, err)
	}
	fmt.Println("\rFile Attachments: ", len(attachments))

	tadpoles_api.DownloadFileAttachments(attachments, backupTarget)
}
