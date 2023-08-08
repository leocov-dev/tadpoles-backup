package api

import (
	"tadpoles-backup/internal/utils"
)

// Response from API
type ParametersResponse struct {
	LastEventTime  utils.JsonTime `json:"last_event_time"`
	FirstEventTime utils.JsonTime `json:"first_event_time"`
	Memberships    []*memberships `json:"memberships"`
}

type memberships struct {
	Dependants []*dependants `json:"dependants"`
}

type dependants struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DisplayName string `json:"display_name"`
	Key         string `json:"person"`
}
