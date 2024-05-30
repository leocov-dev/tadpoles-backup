// This reset command is here to facilitate resetting your password in the
// situation when a tadpoles.com account has been configured to use Google-Auth
// instead of direct login. In this situation it's not possible to reset your
// password via the website or disable the Google-Auth connection. This
// command will call the correct APIs to reset the direct login password.

package commands

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/api"
	"tadpoles-backup/internal/utils"
)

var resetPasswordCmd = &cobra.Command{
	Use:   "reset-password <email-address>",
	Short: "Trigger tadpoles.com to send a password reset email",
	Run:   resetPasswordRun,
	Args:  resetPasswordArgs(),
	PreRun: func(cmd *cobra.Command, args []string) {
		if config.IsNotInteractive() {
			utils.CmdFailed(errors.New("can't run this command in non-interactive mode"))
		}
	},
}

func init() {
	rootCmd.AddCommand(resetPasswordCmd)
}

func resetPasswordArgs() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("<email-address> argument missing")
		}
		return nil
	}
}

func resetPasswordRun(cmd *cobra.Command, args []string) {
	if config.Provider.Value != config.Tadpoles {
	}

	provider := api.GetProvider()

	err := provider.ResetUserPassword(args[0])
	if err != nil {
		utils.CmdFailed(err)
	}
}
