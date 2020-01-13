package commands

import (
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/internal/input"
	"github.com/leocov-dev/tadpoles-backup/internal/tadpoles_api"
	"github.com/spf13/cobra"
)

var (
	statCmd = &cobra.Command{
		Use:   "stat",
		Short: "Print info about events",
		Run:   statRun,
	}
)

func init() {
}

func statRun(cmd *cobra.Command, args []string) {
	fmt.Println("StatRun")
	email, password := input.Credentials()
	fmt.Printf("Email: %s, Password: %s\n", email, password)

	tadpoles_api.PostLogin(email, password)
}
