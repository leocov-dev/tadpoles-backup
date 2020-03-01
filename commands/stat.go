package commands

import (
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/internal/tadpoles_api"
	"github.com/leocov-dev/tadpoles-backup/internal/user_input"
	"github.com/leocov-dev/tadpoles-backup/internal/utils"
	"github.com/spf13/cobra"
)

var (
	statCmd = &cobra.Command{
		Use:   "stat",
		Short: "Print Account Info",
		Run:   statRun,
		PreRun: func(cmd *cobra.Command, args []string) {
			user_input.DoLoginIfNeeded()
		},
	}
)

func statRun(cmd *cobra.Command, args []string) {
	h := utils.NewHeading(":", 15)
	s := utils.StartSpinner("Getting Account Info...")

	info, err := tadpoles_api.GetAccountInfo()
	if err != nil {
		utils.CmdFailed(cmd, err)
	}
	s.Stop()

	h.Write("Timeframe", fmt.Sprintf("%s to %s",
		info.FirstEvent.Format("2006-01-02"),
		info.LastEvent.Format("2006-01-02")))

	h.Write("Children", "")
	for i, dep := range info.Dependants {
		i += 1
		h.WriteAligned(fmt.Sprintf("%d", i), dep, utils.Right)
	}

	s = utils.StartSpinner("Checking Events...")
	attachments, err := tadpoles_api.GetEventAttachmentData(info.FirstEvent, info.LastEvent)
	if err != nil {
		utils.CmdFailed(cmd, err)
	}
	s.Stop()
	h.Write("Pictures/Videos", fmt.Sprintf("%d", len(attachments)))
}
