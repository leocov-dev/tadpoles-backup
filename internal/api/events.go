package api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/client"
	"tadpoles-backup/internal/utils"
	"time"
)

type Attachment struct {
	AttachmentKey string `json:"key"`
	MimeType      string `json:"mime_type"`
}

type Event struct {
	Comment         string         `json:"comment"`
	Attachments     []*Attachment  `json:"new_attachments"`
	ChildName       string         `json:"parent_member_display"`
	EventTime       utils.JsonTime `json:"event_time"`
	CreateTime      utils.JsonTime `json:"create_time"`
	TimeZone        string         `json:"tz"`
	EventKey        string         `json:"key"`
	EventType       string         `json:"type"`
	LocationDisplay string         `json:"location_display"`
}

type Events []*Event

func (e *Event) String() string {
	val, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		log.Error(err)
		return ""
	}
	return string(val)
}

type pageResponse struct {
	Cursor string `json:"cursor"`
	Events Events `json:"events"`
}

func GetEvents(firstEventTime time.Time, lastEventTime time.Time) (events Events, err error) {
	log.Debug(fmt.Sprintf("EventsURL: %s", client.EventsEndpoint))

	params := url.Values{
		"direction":           {"range"},
		"earliest_event_time": {strconv.FormatInt(firstEventTime.Unix(), 10)},
		"latest_event_time":   {strconv.FormatInt(lastEventTime.Unix(), 10)},
		"num_events":          {fmt.Sprint(config.EventsQueryPageSize)},
		"cursor":              nil, // it is acceptable to start cursor as empty
	}

	for true {
		log.Debug("Cursor: ", params.Get("cursor"))
		err = getEventPage(&params, &events)
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

func getEventPage(params *url.Values, events *Events) error {
	urlBase, _ := url.Parse(client.EventsEndpoint)
	urlBase.RawQuery = params.Encode()

	log.Debug("Query: ", urlBase.String())
	resp, err := client.GetApiClient().Get(urlBase.String())
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return client.NewRequestError(resp)
	}

	defer utils.CloseWithLog(resp.Body)
	body, _ := io.ReadAll(resp.Body)

	var page pageResponse
	err = json.Unmarshal(body, &page)

	params.Set("cursor", page.Cursor)

	*events = append(*events, page.Events...)

	return nil
}
