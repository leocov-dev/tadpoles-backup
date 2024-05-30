package tadpoles

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"sort"
	"tadpoles-backup/internal/utils"
	"tadpoles-backup/pkg/async"
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
	Events Events `json:"Events"`
}

func fetchEventsPage(client *http.Client, eventsUrl *url.URL) (newEvents Events, cursor string, err error) {

	log.Debug("Query: ", eventsUrl.String())
	resp, err := client.Get(eventsUrl.String())
	if err != nil {
		return nil, "", err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, "", utils.NewRequestError(resp, "could not get Events page")
	}

	defer utils.CloseWithLog(resp.Body)
	body, _ := io.ReadAll(resp.Body)

	var page pageResponse
	err = json.Unmarshal(body, &page)

	return page.Events, page.Cursor, nil
}

func fetchAllEvents(
	ctx context.Context,
	eventCache *ApiCache,
	httpClient *http.Client,
	ep endpoints,
	firstEventTime time.Time,
	lastEventTime time.Time,
	useCache bool,
) (Events, error) {
	var allEvents Events
	var newEvents Events

	if useCache {
		cachedEvents, err := eventCache.readEventCache()
		if err != nil {
			return nil, err
		}
		allEvents = append(allEvents, cachedEvents...)

		cachedEventsLen := len(cachedEvents)
		if cachedEventsLen > 0 {
			lastCachedEvent := cachedEvents[cachedEventsLen-1]
			lastEventTime = lastCachedEvent.EventTime.Time()
			// TODO this falls apart if `end` is not after `lastEventTime`
			firstEventTime = lastEventTime.Add(1 * time.Second)
		}
	}

	pageNum := 0

	// need a non-empty value to enter the while loop
	cursor := "initialize"

	for cursor != "" {
		select {
		case <-ctx.Done():
			return nil, async.NewCanceledError()
		default:
			log.Debug(fmt.Sprintf("Page: %d Cursor: %s", pageNum, cursor))
			var pageEvents Events
			var pageError error

			pageEvents, cursor, pageError = fetchEventsPage(httpClient, ep.eventsUrl(firstEventTime, lastEventTime, cursor))
			if pageError != nil {
				return nil, pageError
			}
			newEvents = append(newEvents, pageEvents...)
			pageNum += 1
		}
	}

	if len(newEvents) >= 0 {
		log.Debug(fmt.Sprintf("Adding New Events: %d", len(newEvents)))
		if useCache {
			cacheErr := eventCache.updateEventCache(newEvents)
			if cacheErr != nil {
				return nil, cacheErr
			}
		}
		allEvents = append(allEvents, newEvents...)
	}

	log.Debug("Get Events Done...")

	allEvents.Sort(func(e1, e2 *Event) bool {
		return e1.EventTime.Time().Before(e2.EventTime.Time())
	})

	return allEvents, nil
}
