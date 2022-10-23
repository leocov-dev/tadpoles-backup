//go:build windows
// +build windows

package config

import "time"

var (
	SpinnerCharSet = []string{"|", "/", "-", "\\"}
	SpinnerSpeed   = time.Duration(75)
)
