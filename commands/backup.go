package commands

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"os"
	"path/filepath"
	"runtime"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/schemas"
	"tadpoles-backup/internal/tadpoles"
	"tadpoles-backup/internal/user_input"
	"tadpoles-backup/internal/utils"
	"tadpoles-backup/internal/utils/progress"
	"tadpoles-backup/internal/utils/spinners"
)

var (
	backupCmd = &cobra.Command{
		Use:   "backup <target-directory>",
		Short: "Backup New Images.",
		Run:   backupRun,
		Args:  backupArgs(),
		PreRun: func(cmd *cobra.Command, args []string) {
			log.Debug("Backup PersistentPreRun")
			utils.CloseHandlerWithCallback(func() {
				cancelBackup()
			})
			err := user_input.DoLoginIfNeeded()
			if err != nil {
				utils.CmdFailed(err)
			}
		},
	}

	detailedBackupJson bool
	ctx, cancelBackup  = context.WithCancel(context.Background())
	concurrencyLimit   int
	defaultConcurrency = runtime.NumCPU() + (runtime.NumCPU() / 2)
)

func init() {
	backupCmd.Flags().VarP(
		schemas.NewConcurrencyValue(defaultConcurrency, &concurrencyLimit),
		"concurrency",
		"c",
		fmt.Sprintf("The number of simultaneous downloads allowed, 1 - %d.", config.MaxConcurrency),
	)
	backupCmd.Flags().BoolVarP(
		&detailedBackupJson,
		"with-files",
		"w",
		false,
		"JSON output includes detailed list of files (this is a large amount of data).",
	)

	rootCmd.AddCommand(backupCmd)
}

func backupArgs() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("<target-directory> argument missing")
		}
		return nil
	}
}

func backupRun(cmd *cobra.Command, args []string) {
	s := spinners.StartNewSpinner("Getting Account Info...")

	backupTarget := filepath.Clean(args[0])
	log.Debug("Backing up to: ", backupTarget)
	err := os.MkdirAll(backupTarget, os.ModePerm)
	if err != nil {
		s.Stop()
		utils.CmdFailed(err)
	}

	info, err := tadpoles.GetAccountInfo()
	if err != nil {
		s.Stop()
		utils.CmdFailed(err)
	}
	s.Stop()

	s = spinners.StartNewSpinner("Checking Events...")
	fileAttachments, err := tadpoles.GetEventFileAttachmentData(info.FirstEvent, info.LastEvent)
	if err != nil {
		utils.CmdFailed(err)
	}
	s.Stop()

	newAttachments, err := tadpoles.PruneAlreadyDownloaded(fileAttachments, backupTarget)
	if err != nil {
		utils.CmdFailed(err)
	}

	attachmentMap := tadpoles.GroupAttachmentsByType(newAttachments)
	attachmentMap.PrettyPrint("New Attachments")

	var saveErrors []string
	count := len(newAttachments)
	if count > 0 {
		bw := progress.StartNewProgressBar(count, "Downloading")

		saveErrors = tadpoles.DownloadFileAttachments(newAttachments, backupTarget, ctx, concurrencyLimit, bw)

		bw.Stop()

		tadpoles.PrintErrorList(saveErrors)
	}

	backupOutput := schemas.NewBackupOutput(newAttachments, attachmentMap, saveErrors)
	err = backupOutput.JsonPrint(detailedBackupJson)
	if err != nil {
		utils.CmdFailed(err)
	}
}
