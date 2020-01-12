package main

import (
	"github.com/leocov-dev/tadpoles-backup/commands"
)

func main() {

	err := commands.Execute()

	if err != nil {
		println(err)
	}

}
