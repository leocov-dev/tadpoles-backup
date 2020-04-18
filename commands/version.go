package commands

import (
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/config"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		if config.Version == "" {
			config.Version = "v0.0.0-dev"
		}
		fmt.Printf("%s\n", config.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
