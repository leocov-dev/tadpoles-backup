package commands

import (
	"fmt"
	"github.com/gosuri/uiprogress"
	"github.com/leocov-dev/tadpoles-backup/config"
	"github.com/leocov-dev/tadpoles-backup/internal/schemas"
	"github.com/leocov-dev/tadpoles-backup/internal/tadpoles"
	"github.com/leocov-dev/tadpoles-backup/internal/user_input"
	"github.com/leocov-dev/tadpoles-backup/internal/utils"
	"github.com/leocov-dev/tadpoles-backup/internal/utils/spinners"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"os"
	"path/filepath"
	"runtime"
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
				utils.CmdFailed(cmd, err)
			}
		},
	}

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
		utils.CmdFailed(cmd, err)
	}

	info, err := tadpoles.GetAccountInfo()
	if err != nil {
		s.Stop()
		utils.CmdFailed(cmd, err)
	}
	s.Stop()

	s = spinners.StartNewSpinner("Checking Events...")
	attachments, err := tadpoles.GetEventAttachmentData(info.FirstEvent, info.LastEvent)
	if err != nil {
		utils.CmdFailed(cmd, err)
	}
	s.Stop()

	newAttachments, err := tadpoles.PruneAlreadyDownloaded(attachments, backupTarget)
	if err != nil {
		utils.CmdFailed(cmd, err)
	}

	utils.WriteMain("New Attachments", fmt.Sprint(len(newAttachments)))
	typeMap := tadpoles.GroupAttachmentsByType(newAttachments)
	for k, v := range typeMap {
		utils.WriteSub(k, fmt.Sprint(len(v)))
	}

	count := len(newAttachments)
	if count > 0 {
		uiprogress.Start()
		pb := uiprogress.AddBar(count).
			AppendCompleted().
			PrependElapsed().
			PrependFunc(func(b *uiprogress.Bar) string {
				return fmt.Sprintf("Downloading (%d/%d)", b.Current(), count)
			})

		saveErrors, err := tadpoles.DownloadFileAttachments(newAttachments, backupTarget, ctx, concurrencyLimit, pb)
		if err != nil {
			utils.CmdFailed(cmd, err)
		}

		uiprogress.Stop()

		utils.PrintErrorList(saveErrors)
	}

}
