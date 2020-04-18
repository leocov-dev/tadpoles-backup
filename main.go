package main

import (
	"github.com/leocov-dev/tadpoles-backup/commands"
)

func main() {
	commands.Execute()
}

//go:generate go-bindata -pkg bindata -o ./internal/bindata/bindata.go ./utils/dist
