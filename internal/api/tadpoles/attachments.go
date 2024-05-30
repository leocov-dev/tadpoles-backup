package tadpoles

import (
	"fmt"
	"tadpoles-backup/internal/schemas"
)

func NewMediaFileFromEventAttachment(
	event Event,
	attachment Attachment,
	endpoints schemas.TadpolesApiEndpoints,
) schemas.MediaFile {
	return schemas.NewMediaFile(
		event.Comment,
		event.EventTime.Time(),
		fmt.Sprintf("%s_%s", event.FormatTimeStamp(), event.ChildName),
		endpoints.AttachmentsUrl(event.EventKey, attachment.AttachmentKey),
		attachment.MimeType,
	)
}
