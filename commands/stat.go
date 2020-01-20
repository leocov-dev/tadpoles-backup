package commands

import (
	"encoding/json"
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/internal/tadpoles_api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	statCmd = &cobra.Command{
		Use:   "stat",
		Short: "Print Account Info",
		Run:   statRun,
	}
)

func init() {
}

func statRun(cmd *cobra.Command, args []string) {
	fmt.Println("Getting Account Info...")

	info, err := tadpoles_api.GetAccountInfo()
	if err != nil {
		cmdFailed(cmd, err)
	}

	log.Debug(fmt.Sprintf("info: %+v", info))
	paramsString, _ := json.MarshalIndent(info, "", "  ")
	log.Debug(string(paramsString))

	fmt.Printf(
		"Timeframe: %s to %s\n",
		info.FirstEvent.Format("2006-01-02"),
		info.LastEvent.Format("2006-01-02"),
	)

	fmt.Println("Children :")
	for i, dep := range info.Dependants {
		i += 1
		fmt.Printf("%8d : %s\n", i, dep)
	}

	fmt.Print("Checking Events...")
	attachments, err := tadpoles_api.GetFileAttachments(info.FirstEvent, info.LastEvent)
	if err != nil {
		cmdFailed(cmd, err)
	}
	fmt.Println("\rFile Attachments: ", len(attachments))
}
