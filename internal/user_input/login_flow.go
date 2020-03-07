package user_input

import (
	"bufio"
	"fmt"
	"github.com/gookit/color"
	"github.com/leocov-dev/tadpoles-backup/internal/api"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"syscall"
	"time"
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
			fmt.Printf("Login expires  : %s\n\n", expires.In(time.Local).Format("Mon Jan 2 03:04:05 PM"))
			break
		}
		log.Debug("Admit Error: ", err)
		color.Red.Println("Login failed, please try again...")
	}
}

// get username and password from user user_input
func credentials() (string, string) {
	color.Magenta.Println("Input: tadpoles.com login required...")
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Email: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Password: ")
	bytePassword, _ := terminal.ReadPassword(syscall.Stdin)
	password := string(bytePassword)
	fmt.Println()

	return strings.TrimSpace(username), strings.TrimSpace(password)
}
