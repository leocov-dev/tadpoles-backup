package api

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

type TadpolesLogin struct {
	client   *http.Client
	loginUrl *url.URL
}

func newTadpolesLogin(request *http.Client) *TadpolesLogin {
	return &TadpolesLogin{
		client:   request,
		loginUrl: tadpolesUrl.JoinPath("auth", "login"),
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
		return nil, newRequestError(resp, "tadpoles login failed")
	}

	logrus.Debug("Login successful")
	return l.admit()
}

func (l *TadpolesLogin) admit() (expires *time.Time, err error) {
	return admitAndStoreCookie(l.client)
}
