package commands

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/provider_client"
	"tadpoles-backup/internal/schemas"
	"tadpoles-backup/internal/utils"
	"tadpoles-backup/internal/utils/spinners"
	"time"
)

var (
	statCmd = &cobra.Command{
		Use:   "stat",
		Short: "Print Account Info",
		Run:   statRun,
		PreRun: func(cmd *cobra.Command, args []string) {
			utils.CloseHandlerWithCallback(func() {
				cancelStat()
				spinners.SpinnerManager.StopAll()
			})
		},
	}
	detailedStatJson    bool
	statCtx, cancelStat = context.WithCancel(context.Background())
)

func init() {
	statCmd.Flags().BoolVarP(&detailedStatJson, "with-files", "w", false, "JSON output includes detailed list of files (this is a large amount of data).")
	rootCmd.AddCommand(statCmd)
}

func statRun(_ *cobra.Command, _ []string) {
	provider := provider_client.GetProviderClient()

	err := provider.LoginIfNeeded()
	if err != nil {
		utils.CmdFailed(err)
	}

	// ------------------------------------------------------------------------
	s := spinners.StartNewSpinner("Fetching Account Info...")
	info, err := provider.GetAccountInfo()
	if err != nil {
		s.Stop()
		utils.CmdFailed(err)
	}
	s.Stop()

	if config.IsHumanReadable() {
		info.PrettyPrint()
	}

	// ------------------------------------------------------------------------
	s = spinners.StartNewSpinner("Fetching Media Info...")
	mediaFiles, err := provider.GetAllMediaFiles(
		statCtx,
		info.FirstEvent,
		time.Now(),
		provider.ShouldUseCache("stat"),
	)
	if err != nil {
		s.Stop()
		utils.CmdFailed(err)
	}
	s.Stop()

	// ------------------------------------------------------------------------
	if config.IsHumanReadable() {
		mediaFiles.CountByType().PrettyPrint("All Media Files")
	} else {
		statOutput := NewStatOutput(*info, mediaFiles)
		statOutput.Print(detailedStatJson)
	}
}

// StatOutput
// Formatting schema for printing account info
type StatOutput struct {
	Info       schemas.AccountInfo
	MediaFiles schemas.MediaFiles `json:"files,omitempty"`
	Images     int                `json:"imageCount,omitempty"`
	Videos     int                `json:"videoCount,omitempty"`
	Unknown    int                `json:"unknownCount,omitempty"`
}

func NewStatOutput(
	info schemas.AccountInfo,
	files schemas.MediaFiles,
) *StatOutput {
	fileMap := files.CountByType()
	return &StatOutput{
		Info:       info,
		MediaFiles: files,
		Images:     fileMap["Images"],
		Videos:     fileMap["Videos"],
		Unknown:    fileMap["Unknown"],
	}
}

func (so *StatOutput) Print(detailed bool) {
	if !detailed {
		so.MediaFiles = nil
	}

	jsonString, err := json.Marshal(so)
	if err != nil {
		log.Error(err)
	}

	fmt.Println(string(jsonString))
}
