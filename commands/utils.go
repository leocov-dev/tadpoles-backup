package commands

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

func cmdFailed(cmd *cobra.Command, err error) {
	fmt.Printf("%s command failed..\n", cmd.Name())
	logrus.Debug(err)
	os.Exit(1)
}
