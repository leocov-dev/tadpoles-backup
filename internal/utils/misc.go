package utils

import (
	"fmt"
	"github.com/h2non/filetype/matchers"
	"github.com/h2non/filetype/types"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func CloseWithLog(f io.Closer) {
	err := f.Close()
	if err != nil {
		log.Errorf("failed to close file: %s", err)
	}
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func CopyFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}

	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}

	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	outputFile.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
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
