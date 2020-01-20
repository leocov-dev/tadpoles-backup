package main

import (
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/commands"
	log "github.com/sirupsen/logrus"
)

func main() {

	err := commands.Execute()

	if err != nil {
		fmt.Println("Could not execute the command. Try debug mode for more information.")
		log.Debug(err)
	}

}
