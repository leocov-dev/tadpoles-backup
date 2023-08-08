package api

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
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

func (s *Spec) getEventPage(request *http.Client, params *url.Values, events *Events) error {
	urlBase := s.Endpoints.Events
	urlBase.RawQuery = params.Encode()

	log.Debug("Query: ", urlBase.String())
	resp, err := request.Get(urlBase.String())
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return newRequestError(resp, "could not get events page")
	}

	defer utils.CloseWithLog(resp.Body)
	body, _ := io.ReadAll(resp.Body)

	var page pageResponse
	err = json.Unmarshal(body, &page)

	params.Set("cursor", page.Cursor)

	*events = append(*events, page.Events...)

	return nil
}
