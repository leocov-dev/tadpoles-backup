package commands

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
	"tadpoles-backup/internal/api"
	"tadpoles-backup/internal/utils"
)

var (
	resetOptions    = []string{"all", "login", "cache"}
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
		if len(args) < 1 {
			return fmt.Errorf("specify one of [%s]", resetOptsString)
		}

		choice := args[0]

		found := false
		for _, item := range resetOptions {
			if item == choice {
				found = true
			}
		}

		if !found {
			return fmt.Errorf("specify one of [%s]", resetOptsString)
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

	provider := api.GetProvider()

	switch choice {
	case "login":
		err = provider.ClearLoginData()
	case "cache":
		err = provider.ClearCache()
	case "all":
		allErrors := provider.ClearAll()
		var errStrList []string

		for _, e := range allErrors {
			if e == nil {
				continue
			}
			errStrList = append(errStrList, e.Error())
		}
		if errStrList != nil {
			err = errors.New(strings.Join(errStrList, "; "))
		}
	}

	if err != nil {
		utils.CmdFailed(err)
	}
}
