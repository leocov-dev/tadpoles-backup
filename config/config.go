package config

import (
	"fmt"
	"os"
	"path/filepath"
)

var Name = "tadpoles-backup"
var DotName = fmt.Sprintf(".%s", Name)
var Version string
var EventsQueryPageSize = 300
var TempFilePattern = fmt.Sprintf("%s-*", Name)
var MaxConcurrency int64 = 128
var userHomeDir, _ = os.UserHomeDir()
var TempDir = filepath.Join(os.TempDir(), DotName)
var DataDir = filepath.Join(userHomeDir, DotName)
var TadpolesCookieFile = filepath.Join(DataDir, fmt.Sprintf("%s-cookie", DotName))
var TadpolesDatabaseFile = filepath.Join(DataDir, fmt.Sprintf("%s-cache", DotName))

// Helpers
var makeDirs = []string{TempDir, DataDir}

func init() {
	for _, dir := range makeDirs {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to make dir '%s' %s\n", dir, err.Error())
		}
	}
}

func ClearCookiesFile() error {
	return os.Remove(TadpolesCookieFile)
}

func ClearDatabaseFile() error {
	return os.Remove(TadpolesDatabaseFile)
}
