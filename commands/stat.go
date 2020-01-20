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

	params, err := tadpoles_api.GetParameters()
	if err != nil {
		cmdFailed(cmd, err)
	}
	log.Debug(fmt.Sprintf("params: %+v", params))
	paramsString, _ := json.MarshalIndent(params, "", "  ")
	log.Debug(string(paramsString))

	first := params.FirstEvent.Format("2006-01-02")
	last := params.LastEvent.Format("2006-01-02")
	fmt.Printf("Timeframe: %s to %s\n", first, last)
	fmt.Println("Children :")
	for i, dep := range params.Dependants {
		i += 1
		fmt.Printf("%8d : %s\n", i, dep)
	}
}
