package tadpoles_api

import (
	"github.com/leocov-dev/tadpoles-backup/internal/schemas"
	"time"
)

func GetAccountInfo() (info *schemas.Info, err error) {
	parameters, err := ApiParameters()
	if err != nil {
		return nil, err
	}

	return translateParameters(parameters), nil
}

func GetFileAttachments(firstEventTime time.Time, lastEventTime time.Time) (attachments []*schemas.FileAttachment, err error) {
	events, err := ApiEvents(firstEventTime, lastEventTime)
	if err != nil {
		return nil, err
	}

	return flattenAttachments(events), nil
}
