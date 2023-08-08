package api

import (
	"encoding/json"
	"fmt"
	"github.com/corpix/uarand"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/utils"
	"time"
)

type spec struct {
	Endpoints Endpoints
	request   *http.Client
	Login     Login
}

func (s *spec) GetAttachment(eventKey string, attachmentKey string) (resp *http.Response, err error) {
	params := url.Values{
		"obj": {eventKey},
		"key": {attachmentKey},
	}

	urlBase := s.Endpoints.Attachments
	urlBase.RawQuery = params.Encode()

	resp, err = s.request.Get(urlBase.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, newRequestError(resp, "could not get attachment")
	}

	return resp, err
}

func (s *spec) GetEvents(firstEventTime time.Time, lastEventTime time.Time) (events Events, err error) {
	params := url.Values{
		"direction":           {"range"},
		"earliest_event_time": {strconv.FormatInt(firstEventTime.Unix(), 10)},
		"latest_event_time":   {strconv.FormatInt(lastEventTime.Unix(), 10)},
		"num_events":          {fmt.Sprint(config.EventsQueryPageSize)},
		"cursor":              nil, // it is acceptable to start cursor as empty
	}

	for true {
		log.Debug("Cursor: ", params.Get("cursor"))
		err = s.getEventPage(s.request, &params, &events)
		if err != nil {
			log.Debug("Get Page Error: ", err)
			return events, err
		}

		// cursor will be empty when no more pages
		if params.Get("cursor") == "" {
			log.Debug("Get Events Done...")
			break
		}
	}

	return events, nil
}

func (s *spec) GetParameters() (params *ParametersResponse, err error) {
	resp, err := s.request.Get(s.Endpoints.Parameters.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, newRequestError(resp, "could not get parameters")
	}

	defer utils.CloseWithLog(resp.Body)
	body, _ := io.ReadAll(resp.Body)

	err = json.Unmarshal(body, &params)

	return params, err
}

type tadpolesTransport struct{}

func (t *tadpolesTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", uarand.GetRandom())
	return http.DefaultTransport.RoundTrip(req)
}

func newSpec() *spec {
	spec := &spec{
		request: &http.Client{
			Jar:       deserializeCookies(tadpolesUrl),
			Transport: &tadpolesTransport{},
		},
		Endpoints: newEndpoints(),
	}

	switch config.Provider.String() {
	case config.BRIGHT_HORIZONS:
		log.Debug("using Bright Horizons login")
		spec.Login = newBrightHorizonsLogin(spec.request)
	default:
		log.Debug("using Tadpoles login")
		spec.Login = newTadpolesLogin(spec.request)
	}

	return spec
}

var (
	Spec *spec
)

func SetupAPISpec() {
	Spec = newSpec()
}
