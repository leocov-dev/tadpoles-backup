package tadpoles_api

import (
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/internal/client"
)

func GetAttachments() {
	fmt.Printf("AttachmentsURL: %s", client.AttachmentsEndpoint)
}
