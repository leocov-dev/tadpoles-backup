package api

import (
	"fmt"
	"github.com/corpix/uarand"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strconv"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/login"
	"tadpoles-backup/internal/utils"
	"time"
)

type spec struct {
	Endpoints Endpoints
	request   *http.Client
	Login     login.Login
}

func (s *spec) GetAttachment(eventKey string, attachmentKey string) (resp *http.Response, err error) {
	params := url.Values{
		"obj": {eventKey},
		"key": {attachmentKey},
	}

	urlBase := *s.Endpoints.Attachments
	urlBase.RawQuery = params.Encode()

	resp, err = s.request.Get(urlBase.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, utils.NewRequestError(resp, "could not get attachment")
	}

	return resp, err
}

func (s *spec) GetEvents(firstEventTime time.Time, lastEventTime time.Time) (events Events, err error) {
	first := "0"
	if !firstEventTime.IsZero() {
		first = strconv.FormatInt(firstEventTime.Unix(), 10)
	}

	last := strconv.FormatInt(lastEventTime.Unix(), 10)

	params := url.Values{
		"direction":           {"range"},
		"earliest_event_time": {first},
		"latest_event_time":   {last},
		"num_events":          {fmt.Sprint(config.EventsQueryPageSize)},
		"cursor":              nil, // it is acceptable to start cursor as empty
	}

	pageNum := 0

	// need a non-empty value to enter the while loop
	cursor := "initialize"

	for cursor != "" {
		log.Debug(fmt.Sprintf("Page: %d Cursor: %s", pageNum, cursor))
		err = s.appendEventsPage(s.request, &params, &events)
		if err != nil {
			log.Debug("Get Page Error: ", err)
			return events, err
		}
		pageNum += 1
		cursor = params.Get("cursor")
	}
	log.Debug("Get Events Done...")

	events.Sort(func(e1, e2 *Event) bool {
		return e1.EventTime.Time().Before(e2.EventTime.Time())
	})

	return events, nil
}

type randomUserAgentTransport struct{}

func (t *randomUserAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", uarand.GetRandom())
	return http.DefaultTransport.RoundTrip(req)
}

func newSpec() *spec {
	switch config.Provider.String() {
	case config.BRIGHT_HORIZONS:
		log.Debug("using Bright Horizons login")
		endpoints := newEndpoints("https://mybrightday.brighthorizons.com")
		request := &http.Client{
			Jar:       login.DeserializeCookies(endpoints.Root),
			Transport: &randomUserAgentTransport{},
			Timeout:   60 * time.Second,
		}
		return &spec{
			request:   request,
			Endpoints: endpoints,
			Login:     login.NewBrightHorizonsLogin(request),
		}
	default:
		log.Debug("using Tadpoles login")
		endpoints := newEndpoints("https://www.tadpoles.com")
		request := &http.Client{
			Jar:       login.DeserializeCookies(endpoints.Root),
			Transport: &randomUserAgentTransport{},
			Timeout:   60 * time.Second,
		}
		return &spec{
			request:   request,
			Endpoints: endpoints,
			Login:     login.NewTadpolesLogin(request),
		}
	}
}

var (
	Spec *spec
)

func SetupAPISpec() {
	Spec = newSpec()
}
