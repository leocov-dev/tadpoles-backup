package tadpoles_api

import (
	"github.com/leocov-dev/tadpoles-backup/internal/api"
	"github.com/leocov-dev/tadpoles-backup/internal/schemas"
)

func eventsToAttachments(events []*api.Event) (attachments []*schemas.FileAttachment) {
	for _, event := range events {
		for _, eventAttachment := range event.Attachments {
			att := &schemas.FileAttachment{
				Comment:       event.Comment,
				AttachmentKey: eventAttachment.AttachmentKey,
				EventKey:      event.EventKey,
				ChildName:     event.ChildName,
				CreateTime:    event.CreateTime.Time(),
				EventTime:     event.EventTime.Time(),
			}
			attachments = append(attachments, att)
		}
	}
	return attachments
}

func translateParameters(parameters *api.ParametersResponse) *schemas.Info {
	info := &schemas.Info{
		FirstEvent: parameters.FirstEventTime.Time(),
		LastEvent:  parameters.LastEventTime.Time(),
	}

	for _, item := range parameters.Memberships {
		for _, dep := range item.Dependants {
			info.Dependants = append(info.Dependants, dep.DisplayName)
		}
	}

	return info
}
