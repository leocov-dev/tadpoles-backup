package tadpoles

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"sort"
	"tadpoles-backup/internal/utils"
)

type Attachment struct {
	AttachmentKey string `json:"key"`
	MimeType      string `json:"mime_type"`
}

type Event struct {
	Comment         string         `json:"comment"`
	Attachments     []*Attachment  `json:"new_attachments"`
	ChildName       string         `json:"parent_member_display"`
	EventTime       utils.EpocTime `json:"event_time"`
	CreateTime      utils.EpocTime `json:"create_time"`
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

func (e *Event) FormatTimeStamp() string {
	return fmt.Sprintf("%d%02d%02d%02d%02d%02d",
		e.EventTime.Time().Year(),
		e.EventTime.Time().Month(),
		e.EventTime.Time().Day(),
		e.EventTime.Time().Hour(),
		e.EventTime.Time().Minute(),
		e.EventTime.Time().Second(),
	)
}

type By func(e1, e2 *Event) bool

func (e Events) Sort(by By) {
	es := &eventSorter{
		events: e,
		by:     by,
	}

	sort.Sort(es)
}

type eventSorter struct {
	events Events
	by     By
}

func (s *eventSorter) Len() int {
	return len(s.events)
}

func (s *eventSorter) Swap(i, j int) {
	s.events[i], s.events[j] = s.events[j], s.events[i]
}

func (s *eventSorter) Less(i, j int) bool {
	return s.by(s.events[i], s.events[j])
}

type pageResponse struct {
	Cursor string `json:"cursor"`
	Events Events `json:"events"`
}

func fetchEventsPage(client *http.Client, eventsUrl *url.URL) (newEvents Events, cursor string, err error) {

	log.Debug("Query: ", eventsUrl.String())
	resp, err := client.Get(eventsUrl.String())
	if err != nil {
		return nil, "", err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, "", utils.NewRequestError(resp, "could not get events page")
	}

	defer utils.CloseWithLog(resp.Body)
	body, _ := io.ReadAll(resp.Body)

	var page pageResponse
	err = json.Unmarshal(body, &page)

	return page.Events, page.Cursor, nil
}
