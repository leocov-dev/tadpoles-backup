package user_input

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/api"
	"tadpoles-backup/internal/utils"
	"tadpoles-backup/pkg/headings"
	"time"
)

func DoLoginIfNeeded() error {
	if api.S.Login.NeedsLogin() {
		// no valid cookie, do login
		var email string
		var password string

		if config.HasEnvCreds() {
			email = config.EnvUsername
			password = config.EnvPassword
		} else if config.IsInteractive() {
			email, password = credentials()
		} else {
			utils.CmdFailed(
				errors.New("credentials must be supplied from the environment if running in non-interactive mode"),
			)
		}

		expires, loginError := api.S.Login.DoLogin(email, password)

		if loginError != nil {
			if config.IsHumanReadable() {
				utils.WriteError("Login failed", "Please try again...")
			}
			return loginError
		}

		if config.IsHumanReadable() {
			utils.WriteInfo("Login expires", expires.In(time.Local).Format("Mon Jan 2 03:04:05 PM"))
			fmt.Println("")
		}
	}

	return nil
}

// get username and password from user user_input
func credentials() (string, string) {
	utils.WriteInfo("Input", "tadpoles.com login required...")
	reader := bufio.NewReader(os.Stdin)

	utils.WriteInfo("Email", "", headings.NoNewLine)
	username, _ := reader.ReadString('\n')

	utils.WriteInfo("Password", "", headings.NoNewLine)
	bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	password := string(bytePassword)
	fmt.Println()

	return strings.TrimSpace(username), strings.TrimSpace(password)
}
