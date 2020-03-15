package tadpoles

import (
	"github.com/gosuri/uiprogress"
	"github.com/korovkin/limiter"
	"github.com/leocov-dev/tadpoles-backup/internal/api"
	"github.com/leocov-dev/tadpoles-backup/internal/schemas"
	"strings"
	"time"
)

func GetAccountInfo() (info *schemas.Info, err error) {
	parameters, err := api.Parameters()
	if err != nil {
		return nil, err
	}

	return translateParameters(parameters), nil
}

func GetEventAttachmentData(firstEventTime time.Time, lastEventTime time.Time) (attachments []*schemas.FileAttachment, err error) {
	events, err := api.Events(firstEventTime, lastEventTime)
	if err != nil {
		return nil, err
	}

	attachments = eventsToAttachments(events)

	return attachments, nil
}

func DownloadFileAttachments(attachments []*schemas.FileAttachment, backupTarget string, concurrencyLimit int, progressBar *uiprogress.Bar) (int, []string, error) {
	err := checkAlreadyDownloaded(attachments, backupTarget)
	if err != nil {
		return 0, nil, err
	}

	errorChan := make(chan string, len(attachments))

	skipped := 0

	limit := limiter.NewConcurrencyLimiter(concurrencyLimit)

	for _, attachment := range attachments {
		currAtt := attachment

		if currAtt.AlreadyDownloaded {
			if progressBar != nil {
				progressBar.Incr()
			}
			skipped += 1
			continue
		}

		currAtt.SetBackupRoot(backupTarget)
		limit.Execute(func() {
			saveFileAttachment(currAtt, progressBar, errorChan)
		})
	}

	limit.Wait()

	close(errorChan)

	var saveErrors []string
	for s := range errorChan {
		saveErrors = append(saveErrors, s)
	}

	return skipped, saveErrors, nil
}

func GroupAttachmentsByType(attachments []*schemas.FileAttachment) map[string][]*schemas.FileAttachment {
	attachmentTypeMap := make(map[string][]*schemas.FileAttachment)

	for _, attachment := range attachments {
		mimeRoot := strings.Split(attachment.EventMime, "/")[0]
		switch mimeRoot {
		case "image":
			imageArray := attachmentTypeMap["Images"]
			imageArray = append(imageArray, attachment)
			attachmentTypeMap["Images"] = imageArray
		case "video":
			videoArray := attachmentTypeMap["Videos"]
			videoArray = append(videoArray, attachment)
			attachmentTypeMap["Videos"] = videoArray
		default:
			unknownArray := attachmentTypeMap["Unknown"]
			unknownArray = append(unknownArray, attachment)
			attachmentTypeMap["Unknown"] = unknownArray
		}
	}

	return attachmentTypeMap
}
