package commands

import (
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/leocov-dev/tadpoles-backup/internal/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	resetCmd           = &cobra.Command{}
	resetDeprecatedMsg = "Use 'clear' command instead - 'reset' will be removed in next major release."
)

func init() {
	err := copier.Copy(&resetCmd, &clearCmd)
	if err != nil {
		log.Fatalf(`Unexpected error, report to developer.\n%s\n`, err)
	}

	resetCmd.Run = resetRun
	resetCmd.Use = fmt.Sprintf("reset [%s]", resetOptsString)
	resetCmd.Hidden = true
	resetCmd.Short = fmt.Sprintf("DEPRECATED: %s\n\n%s\n", resetDeprecatedMsg, clearCmd.Short)
	rootCmd.AddCommand(resetCmd)
}

func resetRun(cmd *cobra.Command, args []string) {
	utils.WriteInfo("DEPRECATED", resetDeprecatedMsg)
	clearRun(cmd, args)
}
