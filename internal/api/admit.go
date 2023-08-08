package api

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

// admitAndStoreCookie Must call admit endpoint before any other requests to get proper auth cookies set
// this can be used both to get cookies after doing a login or to test that deserialized
// cookies are valid
func admitAndStoreCookie(request *http.Client) (*time.Time, error) {
	logrus.Debug("Admit...")

	zone, _ := time.Now().Zone()

	resp, err := request.PostForm(
		apiV1Root.JoinPath("athome", "admit").String(),
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
