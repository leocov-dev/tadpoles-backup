package commands

import (
	"github.com/spf13/cobra"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/schemas"
	"tadpoles-backup/internal/tadpoles"
	"tadpoles-backup/internal/user_input"
	"tadpoles-backup/internal/utils"
	"tadpoles-backup/internal/utils/spinners"
)

var (
	statCmd = &cobra.Command{
		Use:   "stat",
		Short: "Print Account Info",
		Run:   statRun,
		PreRun: func(cmd *cobra.Command, args []string) {
			utils.CloseHandler()
			err := user_input.DoLoginIfNeeded()
			if err != nil {
				utils.CmdFailed(err)
			}
		},
	}
	detailedStatJson bool
)

func init() {
	statCmd.Flags().BoolVarP(&detailedStatJson, "with-files", "w", false, "JSON output includes detailed list of files (this is a large amount of data).")
	rootCmd.AddCommand(statCmd)
}

func statRun(cmd *cobra.Command, _ []string) {

	// ------------------------------------------------------------------------
	s := spinners.StartNewSpinner("Checking Events...")
	events, err := tadpoles.GetAllEvents()
	if err != nil {
		s.Stop()
		utils.CmdFailed(err)
	}
	s.Stop()

	// ------------------------------------------------------------------------
	info := schemas.NewInfoFromEvents(events)
	if config.IsHumanReadable() {
		info.PrettyPrint()
	}

	// ------------------------------------------------------------------------
	s = spinners.StartNewSpinner("Parsing Attachments...")
	attachments, err := tadpoles.GetEventFileAttachmentData(events)
	if err != nil {
		s.Stop()
		utils.CmdFailed(err)
	}
	s.Stop()

	// ------------------------------------------------------------------------
	attachmentMap := tadpoles.GroupAttachmentsByType(attachments)
	if config.IsHumanReadable() {
		attachmentMap.PrettyPrint("All Attachments")
	} else {
		statOutput := schemas.NewStatOutput(info, attachments, attachmentMap)
		statOutput.Print(detailedStatJson)
	}
}
