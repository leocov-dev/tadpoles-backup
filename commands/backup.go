package commands

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/gosuri/uiprogress"
	"github.com/leocov-dev/tadpoles-backup/config"
	"github.com/leocov-dev/tadpoles-backup/internal/tadpoles"
	"github.com/leocov-dev/tadpoles-backup/internal/user_input"
	"github.com/leocov-dev/tadpoles-backup/internal/utils"
	"github.com/leocov-dev/tadpoles-backup/pkg/headings"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

var (
	backupCmd = &cobra.Command{
		Use:   "backup [target-directory]",
		Short: "Backup New Images.",
		Run:   backupRun,
		Args:  backupArgs(),
		PreRun: func(cmd *cobra.Command, args []string) {
			user_input.DoLoginIfNeeded()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return utils.CleanupTempFiles()
		},
	}

	concurrencyLimit   int
	defaultConcurrency = runtime.NumCPU() + (runtime.NumCPU() / 2)
)

func init() {
	backupCmd.Flags().VarP(
		newConcurrencyValue(defaultConcurrency, &concurrencyLimit),
		"concurrency",
		"c",
		fmt.Sprintf("The number of simultaneous downloads allowed, 1 - %d.", config.MaxConcurrency),
	)
	rootCmd.AddCommand(backupCmd)
}

func backupArgs() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("[target-directory] argument missing")
		}
		return nil
	}
}

func backupRun(cmd *cobra.Command, args []string) {
	hYellow := headings.NewHeading(":", 15, headings.WithColor(color.Bold, color.FgYellow))
	hRed := hYellow.Copy(headings.WithColor(color.Bold, color.FgHiRed))
	hRedRight := hRed.Copy(headings.AlightRight())
	hRight := hYellow.Copy(
		headings.AlightRight(),
		headings.WithColor(color.Bold, color.FgGreen),
	)
	s := utils.StartSpinner("Backup Started...")

	backupTarget := filepath.Clean(args[0])
	log.Debug("Backing up to: ", backupTarget)
	err := os.MkdirAll(backupTarget, os.ModePerm)
	if err != nil {
		s.Stop()
		utils.CmdFailed(cmd, err)
	}

	info, err := tadpoles.GetAccountInfo()
	if err != nil {
		s.Stop()
		utils.CmdFailed(cmd, err)
	}
	s.Stop()

	s = utils.StartSpinner("Checking Events...")
	log.Debug("") // newline for debug mode
	attachments, err := tadpoles.GetEventAttachmentData(info.FirstEvent, info.LastEvent)
	if err != nil {
		utils.CmdFailed(cmd, err)
	}
	s.Stop()

	hYellow.Write("Attachments", fmt.Sprint(len(attachments)))
	typeMap := tadpoles.GroupAttachmentsByType(attachments)
	for k, v := range typeMap {
		hRight.Write(k, fmt.Sprint(len(v)))
	}

	boldCyan := color.New(color.FgHiCyan, color.Bold)

	uiprogress.Start()
	progressBar := uiprogress.AddBar(len(attachments)).
		AppendCompleted().
		PrependElapsed().
		PrependFunc(func(b *uiprogress.Bar) string {
			return boldCyan.Sprint("Downloading")
		})

	skippedCount, saveErrors, err := tadpoles.DownloadFileAttachments(attachments, backupTarget, concurrencyLimit, progressBar)
	if err != nil {
		utils.CmdFailed(cmd, err)
	}

	uiprogress.Stop()

	hYellow.Write("Skipped", fmt.Sprint(skippedCount))

	if saveErrors != nil {
		hRed.Write("Errors", "")
		for i, e := range saveErrors {
			hRedRight.Write(fmt.Sprint(i+1), e)
		}
		fmt.Println("")
	}
}

// custom concurrency flag for validation
type concurrencyValue int

func newConcurrencyValue(val int, p *int) *concurrencyValue {
	*p = val
	return (*concurrencyValue)(p)
}

func (i *concurrencyValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return err
	}

	if v > config.MaxConcurrency || v < 1 {
		return errors.New(fmt.Sprintf("value must be 1 - %d", config.MaxConcurrency))
	}

	*i = concurrencyValue(v)
	return nil
}

func (i *concurrencyValue) Type() string {
	return "int"
}

func (i *concurrencyValue) String() string { return strconv.FormatInt(int64(*i), 10) }
