package bright_horizons

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"tadpoles-backup/internal/interfaces"
	"tadpoles-backup/internal/utils"
)

type ProfileResponse struct {
	UserId string `json:"id"`
}

func fetchProfile(client interfaces.HttpClient, profileUrl *url.URL) (profile *ProfileResponse, err error) {
	resp, err := client.Get(profileUrl.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, utils.NewRequestError(resp, "Failed to fetch bright horizons user profile")
	}

	defer utils.CloseWithLog(resp.Body)
	body, _ := io.ReadAll(resp.Body)

	err = json.Unmarshal(body, &profile)

	return profile, err
}
