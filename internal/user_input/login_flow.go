package user_input

import (
	"bufio"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/internal/api"
	"github.com/leocov-dev/tadpoles-backup/internal/utils"
	"github.com/leocov-dev/tadpoles-backup/pkg/headings"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"net/url"
	"os"
	"strings"
	"time"
)

func DoLoginIfNeeded() error {
	_, err := api.PostAdmit()

	if err == nil {
		// serialized credential cookie was valid!
		return nil
	}

	if err, ok := err.(*url.Error); ok {
		switch e := err.Unwrap().(type) {
		case x509.CertificateInvalidError:
			return errors.New(fmt.Sprintf("cannot connect to tadpoles.com: %s", e.Error()))
		default:
			return err
		}
	}

	log.Debug("Admit Error: ", err)

	for {
		email, password := credentials()
		err := api.PostLogin(email, password)
		if err != nil {
			log.Debug("Login Error: ", err)
		}
		expires, err := api.PostAdmit()
		if err == nil {
			utils.WriteInfo("Login expires", expires.In(time.Local).Format("Mon Jan 2 03:04:05 PM"))
			fmt.Println("")
			break
		}
		log.Debug("Admit Error: ", err)
		utils.WriteError("Login failed", "Please try again...")
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
