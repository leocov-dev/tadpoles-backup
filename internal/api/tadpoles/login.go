package tadpoles

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"tadpoles-backup/internal/api"
	"tadpoles-backup/internal/utils"
	"time"
)

func loginAdmit(client *http.Client, admitUrl *url.URL, cookieFile string) (expires *time.Time, err error) {
	logrus.Debug("Admit...")

	zone, _ := time.Now().Zone()

	resp, err := client.PostForm(
		admitUrl.String(),
		url.Values{
			"tz": {zone},
		},
	)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, utils.NewRequestError(resp, "tadpoles admit failed")
	}

	logrus.Debug("Admit successful")

	return api.SerializeResponseCookies(cookieFile, resp)
}

func login(
	client *http.Client,
	loginUrl *url.URL,
	email, password string,
) error {
	resp, err := client.PostForm(
		loginUrl.String(),
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
		return utils.NewRequestError(resp, "tadpoles login failed")
	}

	return nil
}
