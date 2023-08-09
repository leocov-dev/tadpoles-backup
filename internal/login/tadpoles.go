package login

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"tadpoles-backup/internal/utils"
	"time"
)

type TadpolesLogin struct {
	client   *http.Client
	loginUrl *url.URL
}

func NewTadpolesLogin(request *http.Client) *TadpolesLogin {
	loginUrl, _ := url.Parse("https://www.tadpoles.com/auth/login")
	return &TadpolesLogin{
		client:   request,
		loginUrl: loginUrl,
	}
}

func (l *TadpolesLogin) NeedsLogin() bool {
	_, err := l.admit()

	return err != nil
}

func (l *TadpolesLogin) DoLogin(email string, password string) (*time.Time, error) {
	logrus.Debug("Login...")
	resp, err := l.client.PostForm(
		l.loginUrl.String(),
		url.Values{
			"email":    {email},
			"password": {password},
			"service":  {"tadpoles"},
		},
	)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, utils.NewRequestError(resp, "tadpoles login failed")
	}

	logrus.Debug("Login successful")
	return l.admit()
}

func (l *TadpolesLogin) admit() (expires *time.Time, err error) {
	logrus.Debug("Admit...")

	zone, _ := time.Now().Zone()

	admitUrl, _ := url.Parse("https://www.tadpoles.com/remote/v1/athome/admit")
	resp, err := l.client.PostForm(
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

	return serializeResponseCookies(resp)
}
