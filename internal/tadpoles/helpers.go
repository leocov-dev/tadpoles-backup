package tadpoles

import (
	"github.com/leocov-dev/tadpoles-backup/internal/api"
	"github.com/leocov-dev/tadpoles-backup/internal/schemas"
	log "github.com/sirupsen/logrus"
)

func eventsToAttachments(events []*api.Event) (attachments []*schemas.FileAttachment) {
	for _, event := range events {
		for _, eventAttachment := range event.Attachments {
			// skip pdf files
			if eventAttachment.MimeType == "application/pdf" {
				log.Debugf("skipping pdf: %s@%s \n", event.ChildName, event.EventTime)
				continue
			}
			att := &schemas.FileAttachment{
				Comment:       event.Comment,
				AttachmentKey: eventAttachment.AttachmentKey,
				EventKey:      event.EventKey,
				ChildName:     event.ChildName,
				EventTime:     event.EventTime.Time(),
				EventMime:     eventAttachment.MimeType,
			}
			attachments = append(attachments, att)
		}
	}
	return attachments
}