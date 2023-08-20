package bright_horizons

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"tadpoles-backup/internal/utils"
	"time"
)

type Dependent struct {
	Id          string    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	FirstRecord time.Time `json:"earliest_memory"`
	LastRecord  time.Time `json:"graduation_date"`
}

func (d *Dependent) DisplayName() string {
	return fmt.Sprintf("%s %s", d.FirstName, d.LastName)
}

type Dependents []Dependent

func fetchDependents(client *http.Client, dependentUrl *url.URL) (dependents Dependents, err error) {
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
