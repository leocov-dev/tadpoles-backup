package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"tadpoles-backup/config"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", config.GetVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
