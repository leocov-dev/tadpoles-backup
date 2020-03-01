package commands

import (
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s %s\n", config.Name, config.Version)
	},
}
