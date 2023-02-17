package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"tadpoles-backup/config"
)

type jsonError struct {
	Error string `json:"error"`
}

func CmdFailed(err error) {
	if config.JsonOutput {
		errorInt := jsonError{
			Error: err.Error(),
		}
		jsonString, _ := json.Marshal(errorInt)
		fmt.Println(string(jsonString))
	} else {
		WriteError("Cmd Error", err.Error())
	}
	os.Exit(1)
}
