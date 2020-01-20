package commands

import (
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/internal/user_input"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "tadpoles-backup",
		Short: "Backup photos of your child from www.tadpoles.com",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// These run before every sub-command
			setLoggingLevel()
			user_input.DoLoginIfNeeded()
		},
	}
	debugMode   bool
	verboseMode bool
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "Print additional debug and informational logs. Setting this always sets --verbose as well.")
	rootCmd.PersistentFlags().BoolVarP(&verboseMode, "verbose", "v", false, "Print additional informational logs.")

	rootCmd.AddCommand(statCmd)
	rootCmd.AddCommand(backupCmd)
}

func Execute() error {
	return rootCmd.Execute()
}

func setLoggingLevel() {
	if debugMode {
		log.SetLevel(log.DebugLevel)
		fmt.Println("*** In Debug Mode ***")
		fmt.Println()
	} else if verboseMode {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}
}
