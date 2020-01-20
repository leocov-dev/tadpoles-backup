package tadpoles_api

import (
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/internal/client"
)

func GetEvents() {
	fmt.Printf("EventsURL: %s", client.EventsEndpoint)
}
