package commands

import (
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
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
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			_ = os.RemoveAll(config.TempDir)
		},
	}
	debugMode bool
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "Print additional debug and informational logs.")
	_ = rootCmd.PersistentFlags().MarkHidden("debug")
}

func Execute() {
	_ = rootCmd.Execute()
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
