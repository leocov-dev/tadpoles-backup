package commands

import (
	"encoding/json"
	"fmt"
	"github.com/gosuri/uiprogress"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"os"
	"path/filepath"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/provider_client"
	"tadpoles-backup/internal/schemas"
	"tadpoles-backup/internal/utils"
	"tadpoles-backup/internal/utils/progress"
	"tadpoles-backup/internal/utils/spinners"
	"time"
)

var (
	backupCmd = &cobra.Command{
		Use:   "backup <target-directory>",
		Short: "Backup New Images.",
		Run:   backupRun,
		Args:  backupArgs(),
		PreRun: func(cmd *cobra.Command, args []string) {
			utils.CloseHandlerWithCallback(func() {
				spinners.SpinnerManager.StopAll()
				uiprogress.Stop()
				cancelBackup()
			})
		},
	}

	detailedBackupJson      bool
	backupCtx, cancelBackup = context.WithCancel(context.Background())
)

func init() {
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

func backupRun(_ *cobra.Command, args []string) {
	provider := provider_client.GetProviderClient()

	err := provider.LoginIfNeeded()
	if err != nil {
		utils.CmdFailed(err)
	}

	// ------------------------------------------------------------------------
	backupTarget := filepath.Clean(args[0])
	log.Debug("Backing up to: ", backupTarget)
	err = os.MkdirAll(backupTarget, os.ModePerm)
	if err != nil {
		utils.CmdFailed(err)
	}

	// ------------------------------------------------------------------------
	s := spinners.StartNewSpinner("Fetching Media Data...")
	info, err := provider.GetAccountInfo()
	if err != nil {
		s.Stop()
		utils.CmdFailed(err)
	}
	mediaFiles, err := provider.GetAllMediaFiles(
		backupCtx,
		info.FirstEvent,
		time.Now(),
		provider.ShouldUseCache("backup"),
	)
	if err != nil {
		s.Stop()
		utils.CmdFailed(err)
	}
	s.Stop()

	// ------------------------------------------------------------------------
	newMediaFiles, err := mediaFiles.FilterOnlyNew(backupTarget)
	if err != nil {
		utils.CmdFailed(err)
	}

	if config.IsHumanReadable() {
		newMediaFiles.CountByType().PrettyPrint("New Media Files")
	}

	// ------------------------------------------------------------------------
	count := len(newMediaFiles)
	if count > 0 {
		bw := progress.StartNewProgressBar(count, "Downloading")

		err = newMediaFiles.DownloadAll(
			provider.GetHttpClient(),
			backupTarget,
			backupCtx,
			bw,
		)

		bw.Stop()
	}

	if config.IsHumanReadable() {
		utils.WriteError("Download Errors:", err.Error())
	} else {
		NewBackupOutput(newMediaFiles, err).Print(detailedBackupJson)
	}
}

// BackupOutput
// Formatting schema for printing backup info
type BackupOutput struct {
	MediaFiles schemas.MediaFiles `json:"files,omitempty"`
	Images     int                `json:"imageCount"`
	Videos     int                `json:"videoCount"`
	Unknown    int                `json:"unknownCount"`
	Error      error              `json:"error"`
}

func NewBackupOutput(files schemas.MediaFiles, err error) BackupOutput {
	countMap := files.CountByType()

	return BackupOutput{
		MediaFiles: files,
		Images:     countMap["Images"],
		Videos:     countMap["Videos"],
		Unknown:    countMap["Unknown"],
		Error:      err,
	}
}

func (bo BackupOutput) Print(detailed bool) {
	if !detailed {
		bo.MediaFiles = nil
	}

	jsonString, err := json.Marshal(bo)
	if err != nil {
		log.Error(err)
	}

	fmt.Println(string(jsonString))
}
