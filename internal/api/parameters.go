package api

import (
	"encoding/json"
	"github.com/leocov-dev/tadpoles-backup/internal/client"
	"github.com/leocov-dev/tadpoles-backup/internal/utils"
	"io/ioutil"
	"net/http"
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

func Parameters() (params *ParametersResponse, err error) {
	resp, err := client.ApiClient.Get(client.ParametersEndpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, client.NewRequestError(resp)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &params)

	return params, err
}
