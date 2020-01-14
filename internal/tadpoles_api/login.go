package tadpoles_api

import (
	"encoding/json"
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/config"
	"github.com/leocov-dev/tadpoles-backup/internal/errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

func PostLogin(email string, password string) error {
	resp, err := apiClient.PostForm(
		config.LoginUrl,
		url.Values{
			"email":    {email},
			"password": {password},
			"service":  {"tadpoles"},
		},
	)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf("Error logging into: %s\n", config.TadpolesUrl.Host)
		return &errors.RequestError{Response: resp, Message: "Login failed, please try again."}
	}
	storeAuthCookie(resp)
	return nil
}

func storeAuthCookie(response *http.Response) {
	cookiesData := Jar.Cookies(config.TadpolesUrl)
	jsonString, _ := json.MarshalIndent(cookiesData, "", "  ")

	ioutil.WriteFile("tadpolesSession.json", jsonString, 0644)
}
