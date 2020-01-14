package config

import (
	"fmt"
	"net/url"
)

var (
	TadpolesHost   = "https://www.tadpoles.com"
	TadpolesUrl, _ = url.Parse(TadpolesHost)
	apiUrlV1       = fmt.Sprintf("%s/remote/v1", TadpolesHost)
	EventsUrl      = fmt.Sprintf("%s/events", apiUrlV1)
	AttachmentsUrl = fmt.Sprintf("%s/obj_attachment", apiUrlV1)
	LoginUrl       = fmt.Sprintf("%s/auth/login", TadpolesHost)
)
