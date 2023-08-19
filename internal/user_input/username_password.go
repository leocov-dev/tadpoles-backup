package user_input

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/utils"
	"tadpoles-backup/pkg/headings"
)

func GetUsernameAndPassword() (string, string) {
	var username string
	var password string

	if config.HasEnvCreds() {
		username = config.EnvUsername
		password = config.EnvPassword
	} else if config.IsInteractive() {
		username, password = cliCredentials()
	} else {
		utils.CmdFailed(
			errors.New("credentials must be supplied from the environment if running in non-interactive mode"),
		)
	}

	return username, password
}

// get username and password from user user_input
func cliCredentials() (string, string) {
	utils.WriteInfo(
		"Input",
		fmt.Sprintf("%s login required...", config.Provider.String()),
	)
	reader := bufio.NewReader(os.Stdin)

	utils.WriteInfo("Email", "", headings.NoNewLine)
	username, _ := reader.ReadString('\n')

	utils.WriteInfo("Password", "", headings.NoNewLine)
	bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	password := string(bytePassword)
	fmt.Println()

	return strings.TrimSpace(username), strings.TrimSpace(password)
}
