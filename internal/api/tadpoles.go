package api

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

type TadpolesLogin struct {
	request  *http.Client
	loginUrl *url.URL
	admitUrl *url.URL
}

func newTadpolesLogin(request *http.Client) *TadpolesLogin {
	return &TadpolesLogin{
		request:  request,
		loginUrl: tadpolesUrl.JoinPath("auth", "login"),
		admitUrl: apiV1Root.JoinPath("athome", "admit"),
	}
}

func (l *TadpolesLogin) NeedsLogin() bool {
	_, err := l.admit()

	return err != nil
}

func (l *TadpolesLogin) DoLogin(email string, password string) (*time.Time, error) {
	logrus.Debug("Login...")
	resp, err := l.request.PostForm(
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

// admit Must call admit endpoint before any other requests to get proper auth cookies set
// this can be used both to get cookies after doing a login or to test that deserialized
// cookies are valid
func (l *TadpolesLogin) admit() (expires *time.Time, err error) {
	logrus.Debug("Admit...")

	zone, _ := time.Now().Zone()

	resp, err := l.request.PostForm(
		l.admitUrl.String(),
		url.Values{
			"tz": {zone},
		},
	)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, newRequestError(resp, "tadpoles admit failed")
	}

	serializeResponseCookies(resp)

	inLocalTime := resp.Cookies()[0].Expires

	logrus.Debug("Admit successful")
	return &inLocalTime, nil
}
