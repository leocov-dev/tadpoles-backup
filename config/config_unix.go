//go:build linux || darwin
// +build linux darwin

package config

import "time"

var SpinnerCharSet = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
var SpinnerSpeed = time.Duration(100)
