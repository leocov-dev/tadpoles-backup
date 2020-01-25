package commands

import (
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/internal/tadpoles_api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	backupCmd = &cobra.Command{
		Use:   "backup [target-directory]",
		Short: "Backup New Images.",
		Run:   backupRun,
		Args:  backupArgs(),
	}
)

func backupArgs() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("[target-directory] argument missing")
		}
		return nil
	}
}

func backupRun(cmd *cobra.Command, args []string) {
	fmt.Println("Backup Started ...")

	backupTarget := filepath.Clean(args[0])
	err := os.MkdirAll(backupTarget, os.ModePerm)
	if err != nil {
		cmdFailed(cmd, err)
	}
	log.Debug("Backup to: ", backupTarget)

	info, err := tadpoles_api.GetAccountInfo()
	if err != nil {
		cmdFailed(cmd, err)
	}

	fmt.Print("Checking Events...")
	log.Debug("") // newline for debug mode
	attachments, err := tadpoles_api.GetFileAttachments(info.FirstEvent, info.LastEvent)
	if err != nil {
		cmdFailed(cmd, err)
	}
	fmt.Println("\rFile Attachments: ", len(attachments))

	//for i, attachment := range attachments {
	//	log.Info(fmt.Sprintf("Downloading %6d", i))
	//	data, _ := tadpoles_api.ApiAttachment(attachment)
	//	log.Debug("data len: ", binary.Size(data))
	//
	//	kind, _ := filetype.Match(data)
	//	log.Debug(fmt.Sprintf("Kind: %+v\n", kind.Extension))
	//	known := []string{"png", "jpg"}
	//	if !utils.Contains(known, kind.Extension) {
	//		fmt.Println("extension: ", kind.Extension)
	//	}
	//
	//	contentType := http.DetectContentType(data)
	//	log.Debug("Type: ", contentType)
	//}
}
