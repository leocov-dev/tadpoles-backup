package user_input

import (
	"bufio"
	"crypto/x509"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"net/url"
	"os"
	"strings"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/api"
	"tadpoles-backup/internal/utils"
	"tadpoles-backup/pkg/headings"
	"time"
)

func DoLoginIfNeeded() error {
	_, err := api.PostAdmit()

	if err == nil {
		// serialized credential cookie was valid!
		return nil
	}

	if urlError, ok := err.(*url.Error); ok {
		switch e := urlError.Unwrap().(type) {
		case x509.CertificateInvalidError:
			return errors.New(fmt.Sprintf("cannot connect to tadpoles.com: %s", e.Error()))
		default:
			return urlError
		}
	}

	log.Debug("Admit Error: ", err)

	// no valid cookie, do login
	for {
		var email string
		var password string

		if config.HasEnvCreds() {
			email = config.EnvUsername
			password = config.EnvPassword
		} else if config.IsInteractive() {
			email, password = credentials()
		} else {
			utils.CmdFailed(errors.New("credentials must be supplied from the environment if running in non-interactive mode"))
		}

		loginError := api.PostLogin(email, password)
		if loginError != nil {
			log.Debug("Login Error: ", err)
		}

		expires, admitError := api.PostAdmit()
		if admitError == nil {
			if config.IsHumanReadable() {
				utils.WriteInfo("Login expires", expires.In(time.Local).Format("Mon Jan 2 03:04:05 PM"))
				fmt.Println("")
			}
			break
		}
		log.Debug("Admit Error: ", admitError)
		if config.IsHumanReadable() {
			utils.WriteError("Login failed", "Please try again...")
			if config.IsNotInteractive() {
				return admitError
			}
		} else {
			return admitError
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
