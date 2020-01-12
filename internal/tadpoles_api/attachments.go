package tadpoles_api

import (
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/config"
)

func GetAttachments() {
	fmt.Printf("AttachmentsURL: %s", config.AttachmentsUrl)
}
