package tadpoles_api

import (
	"encoding/json"
	"github.com/leocov-dev/tadpoles-backup/internal/client"
	"io/ioutil"
	"net/http"
)

func GetGuardians() (map[string]interface{}, error) {
	resp, err := client.ApiClient.Get(client.GuardiansEndpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, client.NewRequestError(resp)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	jsonBody := make(map[string]interface{})
	err = json.Unmarshal(body, &jsonBody)

	return jsonBody, err
}
