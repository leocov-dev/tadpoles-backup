package utils

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gosuri/uiprogress"
	"github.com/h2non/filetype/matchers"
	"github.com/h2non/filetype/types"
	"github.com/leocov-dev/tadpoles-backup/config"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

func CloseWithLog(f io.Closer) {
	err := f.Close()
	if err != nil {
		PrintError("failed to close file: %s", err)
	}
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("failed removing original file: %s", err)
	}
	return nil
}

func IsImageType(t types.Type) bool {
	for k := range matchers.Image {
		if k == t {
			return true
		}
	}
	return false
}

func CleanupTempFiles() (err error) {

	tempDir := os.TempDir()

	td, err := os.Open(tempDir)
	if err != nil {
		return err
	}
	defer CloseWithLog(td)

	files, err := td.Readdir(0)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), strings.TrimSuffix(config.TempFilePattern, "*")) {
			err = os.Remove(filepath.Join(tempDir, file.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func PrintError(format string, err error) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	red := color.New(color.FgHiRed).SprintFunc()
	_, _ = fmt.Fprintf(color.Output, format, red(err.Error()))
}

func PrintErrorList(errorMsgs []string) {
	if errorMsgs != nil {
		WriteError("Errors", "")
		for i, e := range errorMsgs {
			WriteErrorSub.Write(fmt.Sprint(i+1), e)
		}
		fmt.Println("")
	}
}

func WithProgressBar(title string, count int, operation func(pb *uiprogress.Bar) []string) {
	uiprogress.Start()
	progressBar := uiprogress.AddBar(count).
		AppendCompleted().
		PrependElapsed()
	if title != "" {
		progressBar.PrependFunc(func(b *uiprogress.Bar) string {
			return title
		})
	}

	errorMsgs := operation(progressBar)

	uiprogress.Stop()

	PrintErrorList(errorMsgs)
}

func CloseHandler() {
	CloseHandlerWithCallback(func() {})
}

func CloseHandlerWithCallback(cb func()) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		time.Sleep(1 * time.Second)
		<-c
		if runtime.GOOS != "windows" {
			// makes the cursor visible
			fmt.Print("\033[?25h")
		}
		fmt.Println("\rCtrl+C pressed in Terminal")
		cb()
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()
}
