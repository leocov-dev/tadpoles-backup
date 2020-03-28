package config

import (
	"fmt"
	"os"
	"path/filepath"
)

var Name = "tadpoles-backup"
var DotName = fmt.Sprintf(".%s", Name)
var Version string
var EventsQueryPageSize = 100
var TempFilePattern = fmt.Sprintf("%s-*", Name)
var MaxConcurrency int64 = 128
var userHomeDir, _ = os.UserHomeDir()
var TempDir = filepath.Join(os.TempDir(), DotName)
var DataDir = filepath.Join(userHomeDir, DotName)
var TadpolesCookieFile = filepath.Join(DataDir, fmt.Sprintf("%s-cookie", DotName))

// Helpers
var makeDirs = []string{TempDir, DataDir}

func init() {
	for _, dir := range makeDirs {
		fmt.Printf("Making: %s\n", dir)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to make dir '%s' %s\n", dir, err.Error())
		}
	}
}
