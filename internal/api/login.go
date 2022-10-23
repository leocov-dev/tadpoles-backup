package api

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"tadpoles-backup/internal/client"
	"time"
)

func PostLogin(email string, password string) error {
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
func PostAdmit() (expires *time.Time, err error) {
	log.Debug("Admit...")
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
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, client.NewRequestError(resp)
	}

	client.SerializeResponseCookies(resp)

	inLocalTime := resp.Cookies()[0].Expires

	log.Debug("Admit successful")
	return &inLocalTime, nil
}
