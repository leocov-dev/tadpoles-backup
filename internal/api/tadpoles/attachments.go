package tadpoles

import (
	"fmt"
	"tadpoles-backup/internal/http_utils"
	"tadpoles-backup/internal/schemas"
)

func NewMediaFileFromEventAttachment(
	event Event,
	attachment Attachment,
	endpoints schemas.TadpolesApiEndpoints,
) http_utils.MediaFile {
	return http_utils.NewMediaFile(
		event.Comment,
		event.EventTime.Time(),
		fmt.Sprintf("%s_%s", event.FormatTimeStamp(), event.ChildName),
		endpoints.AttachmentsUrl(event.EventKey, attachment.AttachmentKey),
		attachment.MimeType,
	)
}
