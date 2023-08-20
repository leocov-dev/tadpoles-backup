package tadpoles

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"tadpoles-backup/internal/api"
	"tadpoles-backup/internal/schemas"
	"time"
)

type ApiSpec struct {
	Client    *http.Client
	endpoints endpoints
}

func NewApiSpec(cookieFile string) *ApiSpec {
	endpoints := newEndpoints()
	return &ApiSpec{
		endpoints: endpoints,
		Client: &http.Client{
			Jar:       api.DeserializeCookies(cookieFile, endpoints.root),
			Transport: &api.RandomUserAgentTransport{},
			Timeout:   60 * time.Second,
		},
	}
}

func (a *ApiSpec) GetAccountParameters() (info *ParametersResponse, err error) {
	return fetchParameters(a.Client, a.endpoints.parametersUrl)
}

func (a *ApiSpec) GetEventMediaFiles(event Event) (schemas.MediaFiles, error) {
	mediaFiles := make(schemas.MediaFiles, len(event.Attachments))

	for i, attachment := range event.Attachments {
		mediaFiles[i] = NewMediaFileFromEventAttachment(event, *attachment, a.endpoints)
	}

	return mediaFiles, nil
}

func (a *ApiSpec) GetEvents(firstEventTime time.Time, lastEventTime time.Time) (events Events, err error) {
	pageNum := 0

	// need a non-empty value to enter the while loop
	cursor := "initialize"

	for cursor != "" {
		log.Debug(fmt.Sprintf("Page: %d Cursor: %s", pageNum, cursor))
		var newEvents Events
		newEvents, cursor, err = fetchEventsPage(a.Client, a.endpoints.eventsUrl(firstEventTime, lastEventTime, cursor))
		if err != nil {
			return nil, err
		}
		events = append(events, newEvents...)
		pageNum += 1
	}
	log.Debug("Get Events Done...")

	events.Sort(func(e1, e2 *Event) bool {
		return e1.EventTime.Time().Before(e2.EventTime.Time())
	})

	return events, nil
}

func (a *ApiSpec) NeedsLogin(cookieFile string) bool {
	_, err := loginAdmit(a.Client, a.endpoints.admitUrl, cookieFile)

	return err != nil
}

func (a *ApiSpec) DoLogin(email string, password string, cookieFile string) (*time.Time, error) {
	log.Debug("Login...")

	err := login(a.Client, a.endpoints.loginUrl, email, password)
	if err != nil {
		return nil, err
	}

	log.Debug("Login successful")
	return loginAdmit(a.Client, a.endpoints.admitUrl, cookieFile)
}
