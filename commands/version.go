package commands

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"tadpoles-backup/config"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run:   versionRun,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func versionRun(cmd *cobra.Command, args []string) {
	version := config.GetVersion()

	if config.JsonOutput {
		versionData := struct {
			Version string `json:"version"`
		}{
			Version: version,
		}
		jsonString, _ := json.Marshal(versionData)
		fmt.Println(string(jsonString))
	} else {
		fmt.Printf("%s\n", version)
	}
}
