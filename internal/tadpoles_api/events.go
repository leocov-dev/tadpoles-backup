package tadpoles_api

import (
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/config"
)

func GetEvents() {
	fmt.Printf("EventsURL: %s", config.EventsUrl)
}
