package commands

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/leocov-dev/tadpoles-backup/internal/bindata"
	"github.com/leocov-dev/tadpoles-backup/internal/utils"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"os"
)

var resetPasswordCmd = &cobra.Command{
	Use:   "reset-password",
	Short: "Reset your tadpoles.com password",
	Run:   resetPasswordRun,
}

func init() {
	rootCmd.AddCommand(resetPasswordCmd)
}

func resetPasswordRun(cmd *cobra.Command, args []string) {
	_, err := fmt.Fprintf(color.Output,
		"%s\n%s\n%s\n%s",
		color.HiMagentaString("** Experimental **"),
		"For details visit: https://github.com/leocov-dev/tadpoles-backup/blob/master/.github/GoogleAccountSignIn.md",
		"Do you want to open a web-browser begin the password reset form?",
		color.HiMagentaString("Press ENTER to continue, Ctrl+C to cancel..."),
	)
	if err != nil {
		utils.CmdFailed(cmd, err)
	}

	reader := bufio.NewReader(os.Stdin)
	_, err = reader.ReadString('\n')
	if err != nil {
		utils.CmdFailed(cmd, err)
	}

	fileData, err := bindata.Asset("utils/dist/reset-tadpoles-password.html")
	if err != nil {
		utils.CmdFailed(cmd, err)
	}
	htmlFile := bytes.NewBuffer(fileData)
	err = browser.OpenReader(htmlFile)
	if err != nil {
		utils.CmdFailed(cmd, err)
	}
}
