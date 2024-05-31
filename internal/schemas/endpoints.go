package schemas

import (
	"net/url"
	"time"
)

type TadpolesApiEndpoints interface {
	AttachmentsUrl(eventKey, attachmentKey string) *url.URL
	EventsUrl(firstEventTime, lastEventTime time.Time, cursor string) *url.URL
}
