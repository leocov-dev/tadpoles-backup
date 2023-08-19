package commands

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"os"
	"path/filepath"
	"runtime"
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
			log.Debug("Backup PersistentPreRun")
			utils.CloseHandlerWithCallback(func() {
				cancelBackup()
			})
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
	mediaFiles, err := provider.GetAllMediaFiles(info.FirstEvent, time.Now())
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
	var saveErrors []string
	count := len(newMediaFiles)
	if count > 0 {
		bw := progress.StartNewProgressBar(count, "Downloading")

		saveErrors = newMediaFiles.DownloadAll(
			provider.GetHttpClient(),
			backupTarget,
			concurrencyLimit,
			ctx,
			bw,
		)

		bw.Stop()
	}

	if config.IsHumanReadable() {
		utils.PrintErrorList(saveErrors)
	} else {
		NewBackupOutput(newMediaFiles, saveErrors).Print(detailedBackupJson)
	}
}

// BackupOutput
// Formatting schema for printing backup info
type BackupOutput struct {
	MediaFiles schemas.MediaFiles `json:"files,omitempty"`
	Images     int                `json:"imageCount"`
	Videos     int                `json:"videoCount"`
	Unknown    int                `json:"unknownCount"`
	Errors     []string           `json:"errors"`
}

func NewBackupOutput(files schemas.MediaFiles, errors []string) BackupOutput {
	countMap := files.CountByType()

	return BackupOutput{
		MediaFiles: files,
		Images:     countMap["Images"],
		Videos:     countMap["Videos"],
		Unknown:    countMap["Unknown"],
		Errors:     errors,
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
