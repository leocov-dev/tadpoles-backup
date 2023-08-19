package tadpoles

import (
	"fmt"
	"tadpoles-backup/internal/schemas"
)

func NewMediaFileFromEventAttachment(
	event Event,
	attachment Attachment,
	endpoints endpoints,
) schemas.MediaFile {
	return schemas.NewMediaFile(
		event.Comment,
		event.EventTime.Time(),
		fmt.Sprintf("%s_%s", event.FormatTimeStamp(), event.ChildName),
		endpoints.attachmentsUrl(event.EventKey, attachment.AttachmentKey),
		attachment.MimeType,
	)
}
