package tadpoles

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"tadpoles-backup/internal/utils"
)

type ParametersResponse struct {
	LastEventTime  utils.EpocTime `json:"last_event_time"`
	FirstEventTime utils.EpocTime `json:"first_event_time"`
	Memberships    []*memberships `json:"memberships"`
}

type memberships struct {
	Dependents []*dependents `json:"dependants"`
}

type dependents struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DisplayName string `json:"display_name"`
	Key         string `json:"person"`
}

func fetchParameters(client *http.Client, paramsUrl *url.URL) (params *ParametersResponse, err error) {
	resp, err := client.Get(paramsUrl.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, utils.NewRequestError(resp, "Failed to fetch tadpoles account parameters")
	}

	defer utils.CloseWithLog(resp.Body)
	body, _ := io.ReadAll(resp.Body)

	err = json.Unmarshal(body, &params)

	return params, err
}
