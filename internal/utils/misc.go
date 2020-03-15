package utils

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/h2non/filetype/matchers"
	"github.com/h2non/filetype/types"
	"github.com/leocov-dev/tadpoles-backup/config"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

func CloseWithLog(f io.Closer) {
	err := f.Close()
	if err != nil {
		fmt.Println(HiRed.Sprintf("failed to close file: %s", err.Error()))
	}
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func StartSpinner(title string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Prefix = fmt.Sprintf("%s ", title)
	err := s.Color("cyan", "bold") // implicit s.Start()
	if err != nil {
		fmt.Println(HiRed.Sprintf("Spinner startup failed: %s", err.Error()))
	}
	return s
}

func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}
	defer CloseWithLog(inputFile)

	outputFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer CloseWithLog(outputFile)

	_, err = io.Copy(outputFile, inputFile)
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

	f, err := os.Open(tempDir)
	if err != nil {
		return err
	}
	defer CloseWithLog(f)

	files, err := f.Readdir(0)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), strings.TrimSuffix(config.TempFilePattern, "*")) {
			err = os.Remove(path.Join(tempDir, file.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
