package api

import (
	"encoding/json"
	"io"
	"net/http"
	"tadpoles-backup/internal/client"
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

func GetParameters() (params *ParametersResponse, err error) {
	resp, err := client.GetApiClient().Get(client.ParametersEndpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, client.NewRequestError(resp)
	}

	defer utils.CloseWithLog(resp.Body)
	body, _ := io.ReadAll(resp.Body)

	err = json.Unmarshal(body, &params)

	return params, err
}
