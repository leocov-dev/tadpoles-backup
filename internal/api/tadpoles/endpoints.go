package tadpoles

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type endpoints struct {
	root          *url.URL
	loginUrl      *url.URL
	apiV1Root     *url.URL
	parametersUrl *url.URL
	admitUrl      *url.URL
	resetUrl      *url.URL
}

func newEndpoints() endpoints {
	rootUrl, _ := url.Parse("https://www.tadpoles.com")
	apiV1Root := rootUrl.JoinPath("remote", "v1")

	return endpoints{
		root:          rootUrl,
		apiV1Root:     apiV1Root,
		loginUrl:      rootUrl.JoinPath("auth", "login"),
		parametersUrl: apiV1Root.JoinPath("parameters"),
		admitUrl:      apiV1Root.JoinPath("athome", "admit"),
		resetUrl:      rootUrl.JoinPath("auth", "forgot"),
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
