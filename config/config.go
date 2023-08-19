package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	TADPOLES        string = "tadpoles"
	BRIGHT_HORIZONS        = "brightHorizons"
)

var (
	// VersionTag
	// this must be exported to set it from build command
	// but should not be accessed directly
	VersionTag string

	DebugMode            = false
	exePath, _           = os.Executable()
	exeDir               = filepath.Dir(exePath)
	Name                 = filepath.Base(exePath)
	DotName              = fmt.Sprintf(".%s", Name)
	MaxConcurrency int64 = 128
	userHomeDir, _       = os.UserHomeDir()
	TempDir              = filepath.Join(os.TempDir(), DotName)
	EnvUsername          = os.Getenv("TADPOLES_USER")
	EnvPassword          = os.Getenv("TADPOLES_PASS")

	allProviders    = []string{TADPOLES, BRIGHT_HORIZONS}
	defaultProvider = TADPOLES
	EnvProvider     = os.Getenv("PROVIDER")
	Provider        = NewProviderConfig(allProviders, defaultProvider)

	NonInteractiveMode bool
	JsonOutput         bool
	dataDir            string
)

func init() {
	makeDirs := []string{TempDir}

	for _, dir := range makeDirs {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to make dir '%s' %s\n", dir, err.Error())
		}
	}
}

func IsInteractive() bool {
	return !NonInteractiveMode
}

func IsNotInteractive() bool {
	return NonInteractiveMode
}

func IsPrintingJson() bool {
	return JsonOutput
}

func IsHumanReadable() bool {
	return !JsonOutput
}

func IsContainerized() bool {
	// check if in docker container
	_, err := os.Stat("/.dockerenv")
	if !os.IsNotExist(err) {
		return true
	}

	// check if in kubernetes
	_, kube := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	if kube {
		return true
	}

	return false
}

func GetDataDir() string {
	if dataDir == "" {
		if IsContainerized() {
			dataDir = filepath.Join(exeDir, DotName)
		} else {
			dataDir = filepath.Join(userHomeDir, DotName)
		}
		_ = os.MkdirAll(dataDir, os.ModePerm)
	}

	return dataDir
}

func HasEnvCreds() bool {
	return EnvUsername != "" && EnvPassword != ""
}

func GetVersion() string {
	if VersionTag != "" {
		return VersionTag
	}

	return "0.0.0-dev"
}
