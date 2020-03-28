// +build windows

package config

import "time"

var SpinnerCharSet = []string{"|", "/", "-", "\\"}
var SpinnerSpeed = time.Duration(75)
