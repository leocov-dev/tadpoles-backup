package bright_horizons

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"tadpoles-backup/internal/utils"
	"time"
)

type DependentData struct {
	DependentId string    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	FirstRecord time.Time `json:"earliest_memory"`
	LastRecord  time.Time `json:"graduation_date"`
}

type DependentResponse []DependentData

func fetchDependents(client *http.Client, dependentUrl *url.URL) (dependents DependentResponse, err error) {
	resp, err := client.Get(dependentUrl.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, utils.NewRequestError(resp, "Failed to fetch dependents data")
	}

	defer utils.CloseWithLog(resp.Body)
	body, _ := io.ReadAll(resp.Body)

	err = json.Unmarshal(body, &dependents)

	return dependents, err
}
