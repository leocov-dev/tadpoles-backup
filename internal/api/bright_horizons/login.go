package bright_horizons

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"tadpoles-backup/internal/utils"
)

func check(client *http.Client, checkUrl *url.URL) bool {
	resp, err := client.Get(checkUrl.String())

	return err != nil || resp.StatusCode != http.StatusOK
}

func login(
	client *http.Client,
	loginUrl *url.URL,
	email, password string,
) (string, error) {
	resp, err := client.PostForm(
		loginUrl.String(),
		url.Values{
			"username": {email},
			"password": {password},
			"response": {"jwt"},
		},
	)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", utils.NewRequestError(resp, "bright horizons login failed")
	}

	defer utils.CloseWithLog(resp.Body)
	body, _ := io.ReadAll(resp.Body)

	return string(body), nil
}

type validateResponse struct {
	ApiKey string `json:"api_key"`
}

func fetchApiKey(
	client *http.Client,
	validateUrl *url.URL,
	token string,
) (apiKey string, err error) {
	logrus.Debug("Validate...")

	resp, err := client.PostForm(
		validateUrl.String(),
		url.Values{
			"token": {token},
		},
	)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", utils.NewRequestError(resp, "bright horizons token validation failed")
	}

	defer utils.CloseWithLog(resp.Body)
	body, _ := io.ReadAll(resp.Body)

	var val validateResponse
	err = json.Unmarshal(body, &val)

	return val.ApiKey, nil
}
