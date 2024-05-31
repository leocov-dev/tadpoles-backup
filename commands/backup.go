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
	"tadpoles-backup/internal/api"
	"tadpoles-backup/internal/http_utils"
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
				cancelBackup()
				spinners.SpinnerManager.StopAll()
				uiprogress.Stop()
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
	provider := api.GetProvider()

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
	info, err := provider.FetchAccountInfo()
	if err != nil {
		s.Stop()
		utils.CmdFailed(err)
	}
	mediaFiles, err := provider.FetchAllMediaFiles(
		backupCtx,
		info.FirstEvent,
		time.Now(),
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
			provider.HttpClient(),
			backupTarget,
			backupCtx,
			bw,
		)
		if err != nil {
			bw.Stop()
			utils.CmdFailed(err)
		}

		bw.Stop()
	}

	if config.IsHumanReadable() {
		fmt.Println("Download complete!")
	} else {
		err = NewBackupOutput(newMediaFiles).Print(detailedBackupJson)
		if err != nil {
			utils.CmdFailed(err)
		}
	}
}

// BackupOutput
// Formatting schema for printing backup info
type BackupOutput struct {
	MediaFiles http_utils.MediaFiles `json:"files,omitempty"`
	Images     int                   `json:"imageCount"`
	Videos     int                   `json:"videoCount"`
	Unknown    int                   `json:"unknownCount"`
}

func NewBackupOutput(files http_utils.MediaFiles) BackupOutput {
	countMap := files.CountByType()

	return BackupOutput{
		MediaFiles: files,
		Images:     countMap["Images"],
		Videos:     countMap["Videos"],
		Unknown:    countMap["Unknown"],
	}
}

func (bo BackupOutput) Print(detailed bool) error {
	if !detailed {
		bo.MediaFiles = nil
	}

	jsonString, err := json.Marshal(bo)
	if err != nil {
		return err
	}

	fmt.Println(string(jsonString))
	return nil
}
