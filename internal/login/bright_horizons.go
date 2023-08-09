package login

import (
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"tadpoles-backup/internal/utils"
	"time"
)

type BrightHorizonsLogin struct {
	client      *http.Client
	loginUrl    *url.URL
	validateUrl *url.URL
}

func NewBrightHorizonsLogin(request *http.Client) *BrightHorizonsLogin {
	loginUrl, _ := url.Parse("https://familyinfocenter.brighthorizons.com/mybrightday/login")
	validateUrl, _ := url.Parse("https://mybrightday.brighthorizons.com/auth/jwt/validate")

	return &BrightHorizonsLogin{
		client:      request,
		loginUrl:    loginUrl,
		validateUrl: validateUrl,
	}
}

func (l *BrightHorizonsLogin) NeedsLogin() bool {
	_, err := l.admit()

	return err != nil
}

func (l *BrightHorizonsLogin) DoLogin(email string, password string) (*time.Time, error) {
	resp, err := l.client.PostForm(
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
		return nil, utils.NewRequestError(resp, "bright horizons login failed")
	}
	defer utils.CloseWithLog(resp.Body)
	body, _ := io.ReadAll(resp.Body)
	token := string(body)

	logrus.Debug("JWT Token: ", token)

	return l.validate(string(body))
}

func (l *BrightHorizonsLogin) validate(token string) (expires *time.Time, err error) {
	logrus.Debug("Validate...")

	resp, err := l.client.PostForm(
		l.validateUrl.String(),
		url.Values{
			"token": {token},
		},
	)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, utils.NewRequestError(resp, "bright horizons token validation failed")
	}

	logrus.Debug("Validate successful")

	return l.admit()
}

func (l *BrightHorizonsLogin) admit() (expires *time.Time, err error) {
	return admitAndStoreCookie(l.client)
}
