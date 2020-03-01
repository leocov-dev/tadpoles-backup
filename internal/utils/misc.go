package utils

import (
	"fmt"
	"github.com/briandowns/spinner"
	"os"
	"time"
)

var SpinnerChars = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}

func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
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
	s.Color("cyan", "bold") // implicit s.Start()
	return s
}
