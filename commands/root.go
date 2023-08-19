package commands

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/utils"
)

var (
	rootCmd = &cobra.Command{
		Use: config.Name,
		Long: fmt.Sprintf("%s %s\nBackup photos of your child from www.tadpoles.com",
			config.Name,
			config.GetVersion()),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if config.EnvProvider != "" {
				providerErr := cmd.Flags().Set("provider", config.EnvProvider)
				if providerErr != nil {
					utils.CmdFailed(fmt.Errorf(
						"environment variable PROVIDER should be one of [%s]",
						strings.Join(config.Provider.Allowed, ", "),
					))
				}
			}

			setLoggingLevel()
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			_ = os.RemoveAll(config.TempDir)
		},
	}
)

func init() {
	rootCmd.PersistentFlags().VarP(config.Provider, "provider", "p", fmt.Sprintf("Service provider [%s]", strings.Join(config.Provider.Allowed, ", ")))

	rootCmd.PersistentFlags().BoolVarP(&config.NonInteractiveMode, "non-interactive", "n", false, "Don't use interactive prompts or show dynamic elements.")

	rootCmd.PersistentFlags().BoolVarP(&config.JsonOutput, "json", "j", false, "Output as JSON.")

	rootCmd.PersistentFlags().BoolVarP(&config.DebugMode, "debug", "d", false, "Print additional debug and informational logs.")
	_ = rootCmd.PersistentFlags().MarkHidden("debug")
}

func Execute() {
	_ = rootCmd.Execute()
}

func setLoggingLevel() {
	if config.DebugMode {
		log.SetLevel(log.DebugLevel)
		fmt.Println("*** In Debug Mode ***")
		fmt.Println()
	} else {
		log.SetLevel(log.WarnLevel)
	}
}
