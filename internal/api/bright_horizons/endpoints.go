package bright_horizons

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type endpoints struct {
	root          *url.URL
	apiV1Root     *url.URL
	apiV2Root     *url.URL
	loginUrl      *url.URL
	token2Url     *url.URL
	rvtUrl        *url.URL
	dependentsUrl *url.URL
	parametersUrl *url.URL
}

func newEndpoints() endpoints {
	loginUrl, _ := url.Parse("https://bhlogin.brighthorizons.com")
	bhApiUrl, _ := url.Parse("https://mbdwgateway.brighthorizons.com/api")
	rootUrl, _ := url.Parse("https://mybrightday.brighthorizons.com")
	apiV1Root := rootUrl.JoinPath("remote", "v1")
	apiV2Root := rootUrl.JoinPath("api", "v2")

	rvtUrl := loginUrl
	rvtUrl.RawQuery = url.Values{
		"benefitid":  {"5"},
		"fstargetid": {"1"},
	}.Encode()

	return endpoints{
		root:          rootUrl,
		apiV1Root:     apiV1Root,
		apiV2Root:     apiV2Root,
		loginUrl:      loginUrl,
		token2Url:     bhApiUrl.JoinPath("account", "token2"),
		rvtUrl:        rvtUrl,
		dependentsUrl: bhApiUrl.JoinPath("home", "mychildren"),
		parametersUrl: apiV1Root.JoinPath("parameters"),
	}
}

func (e endpoints) AttachmentsUrl(eventKey, attachmentKey string) *url.URL {
	attachmentsUrl := e.apiV1Root.JoinPath("obj_attachment")

	attachmentsUrl.RawQuery = url.Values{
		"obj": {eventKey},
		"key": {attachmentKey},
	}.Encode()

	return attachmentsUrl
}

func (e endpoints) EventsUrl(firstEventTime, lastEventTime time.Time, cursor string) *url.URL {
	eventsUrl := e.apiV1Root.JoinPath("events")

	first := "0"
	if !firstEventTime.IsZero() {
		first = strconv.FormatInt(firstEventTime.Unix(), 10)
	}

	last := strconv.FormatInt(lastEventTime.Unix(), 10)

	var cursorVal []string
	if cursor != "" && cursor != "initialize" {
		cursorVal = []string{cursor}
	}

	eventsUrl.RawQuery = url.Values{
		"direction":           {"range"},
		"earliest_event_time": {first},
		"latest_event_time":   {last},
		"num_events":          {fmt.Sprint(75)},
		"cursor":              cursorVal,
	}.Encode()

	return eventsUrl
}

func (e endpoints) jwtRedirectUrl(token string) *url.URL {
	jwtUrl := e.apiV2Root.JoinPath("auth", "jwt", "redirect")

	jwtUrl.RawQuery = url.Values{
		"jwt": {token},
	}.Encode()

	return jwtUrl
}

func (e endpoints) samlUrl(redirect string) *url.URL {
	samlUrl, _ := url.Parse(fmt.Sprintf("%s%s", e.loginUrl, redirect))
	return samlUrl
}
