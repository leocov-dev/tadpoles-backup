package commands

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"reflect"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/utils"
)

var debugCmd = &cobra.Command{
	Use:    "debug",
	Short:  "Print debug information",
	Hidden: true,
	Run:    debugRun,
}

func init() {
	rootCmd.AddCommand(debugCmd)
}

func debugRun(cmd *cobra.Command, args []string) {
	debugData := struct {
		Version         string
		Name            string
		TempDir         string
		HasEnvCreds     bool
		IsContainerized bool
		DataDir         string
		CookieFile      string
		CacheDbFile     string
	}{
		Version:         config.GetVersion(),
		Name:            config.Name,
		TempDir:         config.TempDir,
		HasEnvCreds:     config.HasEnvCreds(),
		IsContainerized: config.IsContainerized(),
		DataDir:         config.GetDataDir(),
	}

	if config.IsPrintingJson() {
		jsonData, err := json.Marshal(debugData)
		if err != nil {
			utils.CmdFailed(err)
		}

		fmt.Println(string(jsonData))
	} else {
		v := reflect.ValueOf(debugData)
		typeOfS := v.Type()

		for i := 0; i < v.NumField(); i++ {
			utils.WriteInfo(typeOfS.Field(i).Name, fmt.Sprintf("%v", v.Field(i).Interface()))
		}
	}
}
