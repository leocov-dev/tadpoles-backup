package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/utils"
)

var (
	resetOptions    = []string{"cookie", "cache"}
	resetOptsString = strings.Join(resetOptions, " | ")

	clearCmd = &cobra.Command{
		Use:   fmt.Sprintf("clear [%s]", resetOptsString),
		Short: "Clear all or indicated local data",
		Args:  clearArgs(),
		Run:   clearRun,
	}
)

func init() {
	rootCmd.AddCommand(clearCmd)
}

func clearArgs() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) >= 1 {
			choice := args[0]

			found := false
			for _, item := range resetOptions {
				if item == choice {
					found = true
				}
			}

			if !found {
				return fmt.Errorf("specify one of [%s] or leave blank to reset all", resetOptsString)
			}

		}
		return nil
	}
}

func clearRun(cmd *cobra.Command, args []string) {
	var err error
	var choice string

	if len(args) > 0 {
		choice = args[0]
	}

	switch choice {
	case "cookie":
		err = config.ClearCookiesFile()
	case "cache":
		err = config.ClearCacheFile()
	default:
		err = config.ClearAll()
	}

	if err != nil {
		utils.CmdFailed(cmd, err)
	}
}
