// This reset command is here to facilitate resetting your password in the
// situation when a tadpoles.com account has been configured to use Google-Auth
// instead of direct login. In this situation it's not possible to reset your
// password via the website or disable the Google-Auth connection. This
// command will call the correct APIs to reset the direct login password.

package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"os"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/bindata"
	"tadpoles-backup/internal/utils"
)

var resetPasswordCmd = &cobra.Command{
	Use:   "reset-password",
	Short: "Reset your tadpoles.com password",
	Run:   resetPasswordRun,
	PreRun: func(cmd *cobra.Command, args []string) {
		if config.IsNotInteractive() {
			utils.CmdFailed(errors.New("can't run this command in non-interactive mode"))
		}
	},
}

func init() {
	rootCmd.AddCommand(resetPasswordCmd)
}

func resetPasswordRun(cmd *cobra.Command, args []string) {
	_, err := fmt.Fprintf(color.Output,
		"%s\n%s\n%s\n%s",
		color.HiMagentaString("** Experimental **"),
		"For details visit: https://tadpoles-backup/blob/main/.github/GoogleAccountSignIn.md",
		"Do you want to open a web-browser begin the password reset form?",
		color.HiMagentaString("Press ENTER to continue, Ctrl+C to cancel..."),
	)
	if err != nil {
		utils.CmdFailed(err)
	}

	reader := bufio.NewReader(os.Stdin)
	_, err = reader.ReadString('\n')
	if err != nil {
		utils.CmdFailed(err)
	}

	fileData, err := bindata.Asset("utils/dist/reset-tadpoles-password.html")
	if err != nil {
		utils.CmdFailed(err)
	}
	htmlFile := bytes.NewBuffer(fileData)
	err = browser.OpenReader(htmlFile)
	if err != nil {
		utils.CmdFailed(err)
	}
}
