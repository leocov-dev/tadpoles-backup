// +build !windows

package config

import "time"

var Name = "tadpoles-backup"
var Version string
var EventsQueryPageSize = 100
var TempFilePattern = "tpbk_*"
var MaxConcurrency int64 = 128
var SpinnerCharSet = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
var SpinnerSpeed = time.Duration(100)
