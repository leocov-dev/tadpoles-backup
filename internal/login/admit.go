package login

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"tadpoles-backup/internal/utils"
	"time"
)

// admitAndStoreCookie Must call admit endpoint before any other requests to get proper auth cookies set
// this can be used both to get cookies after doing a login or to test that deserialized
// cookies are valid
func admitAndStoreCookie(request *http.Client) (*time.Time, error) {
	logrus.Debug("Admit...")

	zone, _ := time.Now().Zone()

	admitUrl, _ := url.Parse("https://www.tadpoles.com/remote/v1/athome/admit")
	resp, err := request.PostForm(
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

	serializeResponseCookies(resp)

	inLocalTime := resp.Cookies()[0].Expires

	logrus.Debug("Admit successful")
	return &inLocalTime, nil
}
