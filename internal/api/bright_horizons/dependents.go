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

type Child struct {
	Id          string    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	FirstRecord time.Time `json:"earliest_memory"`
	LastRecord  time.Time `json:"graduation_date"`
}

func (d *Child) DisplayName() string {
	return fmt.Sprintf("%s %s", d.FirstName, d.LastName)
}

type Children []Child

// TODO: no idea what this response actually looks like....
type MyChildrenResponse struct {
	Children []Child `json:"children"`
}

func fetchDependents(client *http.Client, dependentUrl *url.URL) (Children, error) {
	resp, err := client.Get(dependentUrl.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, utils.NewRequestError(resp, "Failed to fetch dependents data")
	}

	defer utils.CloseWithLog(resp.Body)
	body, _ := io.ReadAll(resp.Body)

	var response MyChildrenResponse

	err = json.Unmarshal(body, &response)

	return response.Children, err
}
