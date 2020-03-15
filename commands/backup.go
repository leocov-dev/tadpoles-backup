package commands

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/leocov-dev/tadpoles-backup/internal/tadpoles"
	"github.com/leocov-dev/tadpoles-backup/internal/user_input"
	"github.com/leocov-dev/tadpoles-backup/internal/utils"
	"github.com/leocov-dev/tadpoles-backup/pkg/headings"
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
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return utils.CleanupTempFiles()
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
	hYellow := headings.NewHeading(":", 15, headings.WithColor(color.Bold, color.FgYellow))
	hRed := headings.NewHeading(":", 15, headings.WithColor(color.Bold, color.FgHiRed))
	s := utils.StartSpinner("Backup Started...")

	backupTarget := filepath.Clean(args[0])
	log.Debug("Backing up to: ", backupTarget)
	err := os.MkdirAll(backupTarget, os.ModePerm)
	if err != nil {
		s.Stop()
		utils.CmdFailed(cmd, err)
	}

	info, err := tadpoles.GetAccountInfo()
	if err != nil {
		s.Stop()
		utils.CmdFailed(cmd, err)
	}
	s.Stop()

	s = utils.StartSpinner("Checking Events...")
	log.Debug("") // newline for debug mode
	attachments, err := tadpoles.GetEventAttachmentData(info.FirstEvent, info.LastEvent)
	if err != nil {
		utils.CmdFailed(cmd, err)
	}
	s.Stop()
	hYellow.Write("Attachments", fmt.Sprint(len(attachments)))

	saveErrors, err := tadpoles.DownloadFileAttachments(attachments, backupTarget)
	if err != nil {
		utils.CmdFailed(cmd, err)
	}
	if saveErrors != nil {
		hRed.Write("Errors", "")
		for i, e := range saveErrors {
			//log.Debug("Error...")
			hRed.Write(fmt.Sprint(i+1), e, headings.AlignRight)
		}
		fmt.Println("")
	}
}
