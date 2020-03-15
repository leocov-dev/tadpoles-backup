package api

import (
	"encoding/json"
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/config"
	"github.com/leocov-dev/tadpoles-backup/internal/client"
	"github.com/leocov-dev/tadpoles-backup/internal/utils"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type pageResponse struct {
	Cursor string   `json:"cursor"`
	Events []*Event `json:"events"`
}

type eventAttachment struct {
	AttachmentKey string `json:"key"`
	MimeType      string `json:"mime_type"`
}

type Event struct {
	Comment     string             `json:"comment"`
	Attachments []*eventAttachment `json:"new_attachments"`
	ChildName   string             `json:"parent_member_display"`
	CreateTime  utils.JsonTime     `json:"create_time"`
	EventTime   utils.JsonTime     `json:"event_time"`
	TimeZone    string             `json:"tz"`
	EventKey    string             `json:"key"`
	Member      string             `json:"member"`
}

func Events(firstEventTime time.Time, lastEventTime time.Time) (events []*Event, err error) {
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

func getEventPage(params *url.Values, attachments *[]*Event) error {
	urlBase, _ := url.Parse(client.EventsEndpoint)
	urlBase.RawQuery = params.Encode()

	log.Debug("Query: ", urlBase.String())
	resp, err := client.ApiClient.Get(urlBase.String())
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return client.NewRequestError(resp)
	}

	defer utils.CloseWithLog(resp.Body)
	body, _ := ioutil.ReadAll(resp.Body)

	var page pageResponse
	err = json.Unmarshal(body, &page)

	params.Set("cursor", page.Cursor)

	*attachments = append(*attachments, page.Events...)

	return nil
}
