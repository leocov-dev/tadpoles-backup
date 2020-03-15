package commands

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/leocov-dev/tadpoles-backup/internal/tadpoles"
	"github.com/leocov-dev/tadpoles-backup/internal/user_input"
	"github.com/leocov-dev/tadpoles-backup/internal/utils"
	"github.com/leocov-dev/tadpoles-backup/pkg/headings"
	"github.com/spf13/cobra"
	"time"
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

func init() {
	rootCmd.AddCommand(statCmd)
}

func statRun(cmd *cobra.Command, _ []string) {
	hLeft := headings.NewHeading(":", 15, headings.WithColor(color.Bold, color.FgYellow))
	hRight := hLeft.Copy(
		headings.AlightRight(),
		headings.WithColor(color.Bold, color.FgGreen),
	)
	s := utils.StartSpinner("Getting Account Info...")

	info, err := tadpoles.GetAccountInfo()
	if err != nil {
		utils.CmdFailed(cmd, err)
	}
	s.Stop()

	hLeft.Write("Time-frame", fmt.Sprintf("%s to %s",
		info.FirstEvent.In(time.Local).Format("2006-01-02"),
		info.LastEvent.In(time.Local).Format("2006-01-02")))

	hLeft.Write("Children", "")
	for i, dep := range info.Dependants {
		i += 1
		hRight.Write(fmt.Sprintf("%d", i), dep)
	}

	s = utils.StartSpinner("Checking Events...")
	attachments, err := tadpoles.GetEventAttachmentData(info.FirstEvent, info.LastEvent)
	if err != nil {
		utils.CmdFailed(cmd, err)
	}
	s.Stop()

	hLeft.Write("Attachments", "")
	typeMap := tadpoles.GroupAttachmentsByType(attachments)
	for k, v := range typeMap {
		hRight.Write(k, fmt.Sprint(len(v)))
	}
}
