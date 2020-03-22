package commands

import (
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/internal/tadpoles"
	"github.com/leocov-dev/tadpoles-backup/internal/user_input"
	"github.com/leocov-dev/tadpoles-backup/internal/utils"
	"github.com/leocov-dev/tadpoles-backup/internal/utils/spinners"
	"github.com/spf13/cobra"
	"time"
)

var (
	statCmd = &cobra.Command{
		Use:   "stat",
		Short: "Print Account Info",
		Run:   statRun,
		PreRun: func(cmd *cobra.Command, args []string) {
			utils.CloseHandler()
			user_input.DoLoginIfNeeded()
		},
	}
)

func init() {
	rootCmd.AddCommand(statCmd)
}

func statRun(cmd *cobra.Command, _ []string) {
	s := spinners.StartNewSpinner("Getting Account Info...")

	info, err := tadpoles.GetAccountInfo()
	if err != nil {
		utils.CmdFailed(cmd, err)
	}
	s.Stop()

	utils.WriteMain("Time-frame", fmt.Sprintf("%s to %s",
		info.FirstEvent.In(time.Local).Format("2006-01-02"),
		info.LastEvent.In(time.Local).Format("2006-01-02")))

	utils.WriteMain("Children", "")
	for i, dep := range info.Dependants {
		i += 1
		utils.WriteSub(fmt.Sprintf("%d", i), dep)
	}

	s = spinners.StartNewSpinner("Checking Events...")
	attachments, err := tadpoles.GetEventAttachmentData(info.FirstEvent, info.LastEvent)
	if err != nil {
		utils.CmdFailed(cmd, err)
	}
	s.Stop()

	utils.WriteMain("All Attachments", "")
	typeMap := tadpoles.GroupAttachmentsByType(attachments)
	for k, v := range typeMap {
		utils.WriteSub(k, fmt.Sprint(len(v)))
	}
}
