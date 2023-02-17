package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	Version string

	exePath, _                = os.Executable()
	exeDir                    = filepath.Dir(exePath)
	Name                      = filepath.Base(exePath)
	DotName                   = fmt.Sprintf(".%s", Name)
	EventsQueryPageSize       = 300
	TempFilePattern           = fmt.Sprintf("%s-*", Name)
	MaxConcurrency      int64 = 128
	userHomeDir, _            = os.UserHomeDir()
	TempDir                   = filepath.Join(os.TempDir(), DotName)
	NonInteractiveMode        = false
	JsonOutput                = false
	EnvUsername               = os.Getenv("TADPOLES_USER")
	EnvPassword               = os.Getenv("TADPOLES_PASS")

	dataDir    string
	cookieFile string
	cacheFile  string
)

func init() {
	makeDirs := []string{TempDir, GetDataDir()}

	for _, dir := range makeDirs {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to make dir '%s' %s\n", dir, err.Error())
		}
	}
}

func HasEnvCreds() bool {
	return EnvUsername != "" && EnvPassword != ""
}

func IsContainerized() bool {
	// check if in docker container
	info, err := os.Stat("/.dockerenv")
	if !os.IsNotExist(err) || !info.IsDir() {
		return true
	}

	// check if in kubernetes
	_, kube := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	if kube {
		return true
	}

	return false
}

func GetVersion() string {
	if Version != "" {
		return Version
	}

	return "0.0.0-dev"
}

func GetDataDir() string {
	if dataDir == "" {
		if IsContainerized() {
			dataDir = exeDir
		} else {
			dataDir = filepath.Join(userHomeDir, DotName)
		}
	}
	return dataDir
}

func GetTadpolesCookieFile() string {
	if cookieFile == "" {
		cookieFile = filepath.Join(GetDataDir(), fmt.Sprintf("%s-cookie", DotName))
	}
	return cookieFile
}

func GetCacheDbFile() string {
	if cacheFile == "" {
		cacheFile = filepath.Join(GetDataDir(), fmt.Sprintf("%s-cache", DotName))
	}
	return cacheFile
}

func ClearCookiesFile() error {
	return os.Remove(GetTadpolesCookieFile())
}

func ClearCacheFile() error {
	return os.Remove(GetCacheDbFile())
}

func ClearAll() error {
	cookieErr := ClearCookiesFile()
	cacheErr := ClearCacheFile()

	allErrors := []error{cookieErr, cacheErr}

	var errStr []string

	for _, e := range allErrors {
		if e == nil {
			continue
		}
		errStr = append(errStr, e.Error())
	}

	if errStr != nil {
		return errors.New(strings.Join(errStr, "; "))
	}

	return nil
}
