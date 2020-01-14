package input

import (
	"bufio"
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/internal/tadpoles_api"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"syscall"
)

func DoLoginIfNeeded() {
	for {
		email, password := credentials()
		err := tadpoles_api.PostLogin(email, password)

		if err == nil {
			break
		}
	}
}

func credentials() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Email: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)
	fmt.Println("")

	return strings.TrimSpace(username), strings.TrimSpace(password)
}
