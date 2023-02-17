package schemas

import (
	"encoding/json"
	"fmt"
)

type StatOutput struct {
	*Info
	FileAttachments `json:"files,omitempty"`
	Images          int `json:"imageCount"`
	Videos          int `json:"videoCount"`
	Unknown         int `json:"unknownCount"`
}

func NewStatOutput(info *Info, files FileAttachments, fileMap FileAttachmentMap) StatOutput {
	return StatOutput{
		Info:            info,
		FileAttachments: files,
		Images:          len(fileMap["Images"]),
		Videos:          len(fileMap["Videos"]),
		Unknown:         len(fileMap["Unknown"]),
	}
}

func (so StatOutput) JsonPrint(detailed bool) error {
	if !detailed {
		so.FileAttachments = nil
	}

	jsonString, err := json.Marshal(so)

	if err != nil {
		return err
	}

	fmt.Println(string(jsonString))

	return nil
}
