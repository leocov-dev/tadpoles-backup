package api

import (
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"tadpoles-backup/internal/utils"
	"time"
)

type BrightHorizonsLogin struct {
	request     *http.Client
	loginUrl    *url.URL
	validateUrl *url.URL
	admitUrl    *url.URL
}

func newBrightHorizonsLogin(request *http.Client) *BrightHorizonsLogin {
	loginUrl, _ := url.Parse("https://familyinfocenter.brighthorizons.com/mybrightday/login")

	return &BrightHorizonsLogin{
		request:     request,
		loginUrl:    loginUrl,
		validateUrl: tadpolesUrl.JoinPath("auth", "jwt", "validate"),
		admitUrl:    apiV1Root.JoinPath("athome", "admit"),
	}
}

func (l *BrightHorizonsLogin) NeedsLogin() bool {
	zone, _ := time.Now().Zone()

	resp, err := l.request.PostForm(
		l.admitUrl.String(),
		url.Values{
			"tz": {zone},
		},
	)

	return err != nil || resp.StatusCode != http.StatusOK
}

func (l *BrightHorizonsLogin) DoLogin(email string, password string) (*time.Time, error) {
	resp, err := l.request.PostForm(
		l.loginUrl.String(),
		url.Values{
			"username": {email},
			"password": {password},
			"response": {"jwt"},
		},
	)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, newRequestError(resp, "bright horizons login failed")
	}
	defer utils.CloseWithLog(resp.Body)
	body, _ := io.ReadAll(resp.Body)
	token := string(body)

	// TODO: remove this log
	logrus.Debug("Token:", token)

	return l.validate(token)
}

func (l *BrightHorizonsLogin) validate(token string) (expires *time.Time, err error) {
	logrus.Debug("Validate...")

	resp, err := l.request.PostForm(
		l.validateUrl.String(),
		url.Values{
			"token": {token},
		},
	)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, newRequestError(resp, "bright horizons token validation failed")
	}

	serializeResponseCookies(resp)

	inLocalTime := resp.Cookies()[0].Expires

	logrus.Debug("Validate successful")
	return &inLocalTime, nil
}
