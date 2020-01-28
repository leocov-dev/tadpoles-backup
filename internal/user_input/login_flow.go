package user_input

import (
	"bufio"
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/internal/api"
	"github.com/leocov-dev/tadpoles-backup/internal/client"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"syscall"
)

func DoLoginIfNeeded() {
	client.DeserializeCookies()
	err := api.Admit()

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
		err = api.Admit()
		if err != nil {
			log.Debug("Admit Error: ", err)
		} else {
			// login was successful!
			client.SerializeCookies()
			break
		}

		fmt.Println("Login failed, please try again...")
	}
}

// get username and password from user user_input
func credentials() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Email: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)
	fmt.Println()

	return strings.TrimSpace(username), strings.TrimSpace(password)
}
