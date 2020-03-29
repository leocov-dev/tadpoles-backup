package utils

import (
	"github.com/spf13/cobra"
	"os"
)

func CmdFailed(cmd *cobra.Command, err error) {
	WriteError("Cmd Error", err.Error())
	os.Exit(1)
}
