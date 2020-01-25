package user_input

import (
	"bufio"
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/internal/client"
	"github.com/leocov-dev/tadpoles-backup/internal/tadpoles_api"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"syscall"
)

func DoLoginIfNeeded() {
	client.DeserializeCookies()
	err := tadpoles_api.ApiAdmit()

	if err == nil {
		// serialized credential cookie was valid!
		return
	}

	log.Debug("ApiAdmit Error: ", err)

	for {
		email, password := credentials()
		err := tadpoles_api.ApiLogin(email, password)
		if err != nil {
			log.Debug("Login Error: ", err)
		}
		err = tadpoles_api.ApiAdmit()
		if err != nil {
			log.Debug("ApiAdmit Error: ", err)
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
