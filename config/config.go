package config

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	exePath, _          = os.Executable()
	Name                = filepath.Base(exePath)
	DotName             = fmt.Sprintf(".%s", Name)
	Version             string
	EventsQueryPageSize       = 300
	TempFilePattern           = fmt.Sprintf("%s-*", Name)
	MaxConcurrency      int64 = 128
	userHomeDir, _            = os.UserHomeDir()
	TempDir                   = filepath.Join(os.TempDir(), DotName)
	DataDir                   = filepath.Join(userHomeDir, DotName)
	TadpolesCookieFile        = filepath.Join(DataDir, fmt.Sprintf("%s-cookie", DotName))
	TadpolesCacheFile         = filepath.Join(DataDir, fmt.Sprintf("%s-cache", DotName))
	NonInteractiveMode        = false
	JsonOutput                = false

	// Helpers
	makeDirs = []string{TempDir, DataDir}
)

func GetVersion() string {
	if Version != "" {
		return Version
	}

	return "0.0.0-dev"
}

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

func ClearCacheFile() error {
	return os.Remove(TadpolesCacheFile)
}

func ClearAll() error {
	all := []string{
		TadpolesCookieFile,
		TadpolesCacheFile,
	}
	for _, item := range all {
		err := os.Remove(item)
		if err != nil {
			return err
		}
	}
	return nil
}
