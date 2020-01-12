package config

import "fmt"

var (
	apiUrl         = "https://www.tadpoles.com/remote/v1"
	EventsUrl      = fmt.Sprintf("%s/events", apiUrl)
	AttachmentsUrl = fmt.Sprintf("%s/obj_attachment", apiUrl)
)
