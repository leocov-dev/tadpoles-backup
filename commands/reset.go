package commands

import (
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/config"
	"github.com/leocov-dev/tadpoles-backup/internal/utils"
	"github.com/spf13/cobra"
	"strings"
)

var (
	resetOptions    = []string{"cookie", "cache"}
	resetOptsString = strings.Join(resetOptions, " | ")

	resetCmd = &cobra.Command{
		Use:   fmt.Sprintf("reset [%s]", resetOptsString),
		Short: "Reset all or indicated local data",
		Args:  resetArgs(),
		Run:   resetRun,
	}
)

func init() {
	rootCmd.AddCommand(resetCmd)
}

func resetArgs() cobra.PositionalArgs {
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

func resetRun(cmd *cobra.Command, args []string) {
	var err error
	var choice string

	if len(args) > 0 {
		choice = args[0]
	}

	switch choice {
	case "cookie":
		err = config.ClearCookiesFile()
	case "cache":
		err = config.ClearDatabaseFile()
	default:
		err = config.ClearDatabaseFile()
		if err != nil {
			break
		}
		err = config.ClearCookiesFile()
	}

	if err != nil {
		utils.CmdFailed(cmd, err)
	}
}
