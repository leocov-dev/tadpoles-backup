package schemas

import (
	"encoding/json"
	"fmt"
)

type BackupOutput struct {
	FileAttachments `json:"files,omitempty"`
	Images          int      `json:"imageCount"`
	Videos          int      `json:"videoCount"`
	Unknown         int      `json:"unknownCount"`
	Errors          []string `json:"errors,omitempty"`
}

func NewBackupOutput(files FileAttachments, fileMap FileAttachmentMap) BackupOutput {
	return BackupOutput{
		FileAttachments: files,
		Images:          len(fileMap["Images"]),
		Videos:          len(fileMap["Videos"]),
		Unknown:         len(fileMap["Unknown"]),
	}
}

func (bo BackupOutput) JsonPrint(detailed bool) error {
	if !detailed {
		bo.FileAttachments = nil
	}

	jsonString, err := json.Marshal(bo)

	if err != nil {
		return err
	}

	fmt.Println(string(jsonString))

	return nil
}
