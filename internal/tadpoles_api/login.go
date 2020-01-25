package tadpoles_api

import (
	"github.com/leocov-dev/tadpoles-backup/internal/client"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

func ApiLogin(email string, password string) error {
	log.Debug("Login...")
	resp, err := client.ApiClient.PostForm(
		client.LoginEndpoint,
		url.Values{
			"email":    {email},
			"password": {password},
			"service":  {"tadpoles"},
		},
	)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return client.NewRequestError(resp)
	}

	log.Debug("Login successful")
	return nil
}

// Must call admit endpoint before any other requests to get proper auth cookies set
func ApiAdmit() error {
	log.Debug("ApiAdmit...")
	t := time.Now()
	zone, _ := t.Zone()
	log.Debug("zone: ", zone)
	resp, err := client.ApiClient.PostForm(
		client.AdmitEndpoint,
		url.Values{
			"tz": {zone},
		},
	)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return client.NewRequestError(resp)
	}

	log.Debug("ApiAdmit successful")
	return nil
}
