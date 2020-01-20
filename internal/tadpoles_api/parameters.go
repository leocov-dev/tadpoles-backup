package tadpoles_api

import (
	"encoding/json"
	"github.com/leocov-dev/tadpoles-backup/internal/client"
	"github.com/leocov-dev/tadpoles-backup/internal/utils"
	"io/ioutil"
	"net/http"
	"time"
)

type Parameters struct {
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
}

type Params struct {
	FirstEvent time.Time
	LastEvent  time.Time
	Dependants []string
}

func GetParameters() (*Params, error) {
	resp, err := client.ApiClient.Get(client.ParametersEndpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, client.NewRequestError(resp)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var param Parameters
	err = json.Unmarshal(body, &param)

	var depts []string
	for _, item := range param.Memberships {
		for _, dep := range item.Dependants {
			depts = append(depts, dep.DisplayName)
		}
	}

	params := Params{
		FirstEvent: param.FirstEventTime.Time(),
		LastEvent:  param.LastEventTime.Time(),
		Dependants: depts,
	}

	return &params, err
}
