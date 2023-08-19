package bright_horizons

import (
	"net/url"
	"time"
)

type endpoints struct {
	root        *url.URL
	apiV2Root   *url.URL
	loginUrl    *url.URL
	validateUrl *url.URL
	profileUrl  *url.URL
}

func newEndpoints() endpoints {
	loginUrl, _ := url.Parse("https://familyinfocenter.brighthorizons.com/mybrightday/login")
	rootUrl, _ := url.Parse("https://mybrightday.brighthorizons.com")
	apiV2Root := rootUrl.JoinPath("api", "v2")

	return endpoints{
		root:        rootUrl,
		apiV2Root:   apiV2Root,
		loginUrl:    loginUrl,
		validateUrl: apiV2Root.JoinPath("jwt", "validate"),
		profileUrl:  apiV2Root.JoinPath("user", "profile"),
	}
}

func (e endpoints) dependentsUrl(userId string) *url.URL {
	return e.apiV2Root.JoinPath("dependents", "guardian", userId)
}

func (e endpoints) eventsUrl(childId string, start, end time.Time) *url.URL {
	eventsUrl := e.apiV2Root.JoinPath("dependent", childId, "daily_reports")

	eventsUrl.RawQuery = url.Values{
		"start": {start.Format("2006-01-02")},
		"end":   {end.Format("2006-01-02")},
	}.Encode()

	return eventsUrl
}
