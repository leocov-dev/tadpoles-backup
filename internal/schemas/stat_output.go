package schemas

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"tadpoles-backup/config"
)

type StatOutput struct {
	Info            *Info
	FileAttachments FileAttachments `json:"files,omitempty"`
	Images          int             `json:"imageCount,omitempty"`
	Videos          int             `json:"videoCount,omitempty"`
	Unknown         int             `json:"unknownCount,omitempty"`
}

func NewStatOutput(info *Info, files FileAttachments, fileMap FileAttachmentMap) *StatOutput {
	return &StatOutput{
		Info:            info,
		FileAttachments: files,
		Images:          len(fileMap["Images"]),
		Videos:          len(fileMap["Videos"]),
		Unknown:         len(fileMap["Unknown"]),
	}
}

func (so *StatOutput) Print(detailed bool) {
	if !detailed {
		so.FileAttachments = nil
	}

	if !config.DebugMode {
		so.Unknown = 0
	}

	jsonString, err := json.Marshal(so)

	if err != nil {
		log.Error(err)
	}

	fmt.Println(string(jsonString))
}
