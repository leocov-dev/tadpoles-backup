package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"tadpoles-backup/config"
)

func CmdFailed(err error) {
	if config.IsPrintingJson() {
		errorData := struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		}
		jsonString, _ := json.Marshal(errorData)
		fmt.Println(string(jsonString))
	} else {
		WriteError("Cmd Error", err.Error())
	}
	os.Exit(1)
}
