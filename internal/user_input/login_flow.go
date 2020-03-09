package user_input

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/leocov-dev/tadpoles-backup/internal/api"
	"github.com/leocov-dev/tadpoles-backup/pkg/headings"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"time"
)

var (
	defaultColorOpt = headings.WithColor(color.Bold, color.FgHiMagenta)
	heading         = headings.NewHeading(":", 15, defaultColorOpt)
)

func DoLoginIfNeeded() {
	_, err := api.Admit()

	if err == nil {
		// serialized credential cookie was valid!
		return
	}

	log.Debug("Admit Error: ", err)

	for {
		email, password := credentials()
		err := api.Login(email, password)
		if err != nil {
			log.Debug("Login Error: ", err)
		}
		expires, err := api.Admit()
		if err == nil {
			heading.Write("Login expires", expires.In(time.Local).Format("Mon Jan 2 03:04:05 PM"))
			fmt.Println("")
			break
		}
		log.Debug("Admit Error: ", err)
		color.Red("Login failed, please try again...")
	}
}

// get username and password from user user_input
func credentials() (string, string) {
	heading.Write("Input", "tadpoles.com login required...")
	reader := bufio.NewReader(os.Stdin)

	heading.Write("Email", "", headings.NoNewLine)
	username, _ := reader.ReadString('\n')

	heading.Write("Password", "", headings.NoNewLine)
	bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	password := string(bytePassword)
	fmt.Println()

	return strings.TrimSpace(username), strings.TrimSpace(password)
}
