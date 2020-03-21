package commands

import (
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/config"
	"github.com/leocov-dev/tadpoles-backup/internal/utils"
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
	debugMode bool
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "Print additional debug and informational logs.")
	rootCmd.PersistentFlags().MarkHidden("debug")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		utils.WriteError("Fatal", "Failed to start root command...")
	}
}

func setLoggingLevel() {
	if debugMode {
		log.SetLevel(log.DebugLevel)
		fmt.Println("*** In Debug Mode ***")
		fmt.Println()
	} else {
		log.SetLevel(log.WarnLevel)
	}
}
