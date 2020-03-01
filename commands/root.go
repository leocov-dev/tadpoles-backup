package commands

import (
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: config.Name,
		Long: fmt.Sprintf("%s %s\nBackup photos of your child from www.tadpoles.com",
			config.Name,
			config.Version),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setLoggingLevel()
		},
	}
	debugMode   bool
	verboseMode bool
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "Print additional debug and informational logs. Setting this always sets --verbose as well.")
	rootCmd.PersistentFlags().BoolVarP(&verboseMode, "verbose", "v", false, "Print additional informational logs.")
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
