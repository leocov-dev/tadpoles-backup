package tadpoles

import (
	"context"
	"github.com/gosuri/uiprogress"
	"github.com/leocov-dev/tadpoles-backup/internal/api"
	"github.com/leocov-dev/tadpoles-backup/internal/schemas"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GetAccountInfo() (info *schemas.Info, err error) {
	parameters, err := api.GetParameters()
	if err != nil {
		return nil, err
	}

	return schemas.NewInfoFromParams(parameters), nil
}

func GetEventAttachmentData(firstEventTime time.Time, lastEventTime time.Time) (attachments []*schemas.FileAttachment, err error) {
	events, err := api.GetEvents(firstEventTime, lastEventTime)
	if err != nil {
		return nil, err
	}

	attachments = eventsToAttachments(events)

	return attachments, nil
}

func DownloadFileAttachments(newAttachments []*schemas.FileAttachment, backupRoot string, ctx context.Context, concurrencyLimit int, progressBar *uiprogress.Bar) ([]string, error) {
	errorChan := make(chan string)

	downloadPool := schemas.NewDownloadPool(concurrencyLimit)

	for _, attachment := range newAttachments {
		proc := schemas.NewAttachmentProc(attachment, backupRoot, errorChan, ctx, progressBar)
		downloadPool.Add(proc)
	}
	downloadPool.Process()

	close(errorChan)

	var saveErrors []string
	for s := range errorChan {
		saveErrors = append(saveErrors, s)
	}

	return saveErrors, nil
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

func PruneAlreadyDownloaded(attachments []*schemas.FileAttachment, backupTarget string) (newAttachments []*schemas.FileAttachment, err error) {
	attachmentNames := make(map[string]*schemas.FileAttachment)
	for _, att := range attachments {
		attachmentNames[att.SaveName()] = att
	}

	err = filepath.Walk(backupTarget,
		func(path_ string, info_ os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info_.IsDir() {
				return nil
			}

			file := filepath.Base(path_)
			minusExtension := strings.TrimSuffix(file, filepath.Ext(file))
			logrus.Debugf("minusExtension: %s\n", minusExtension)

			if _, ok := attachmentNames[minusExtension]; ok {
				delete(attachmentNames, minusExtension)
			}

			return nil
		})

	for _, v := range attachmentNames {
		newAttachments = append(newAttachments, v)
	}

	return newAttachments, err
}

func eventsToAttachments(events []*api.Event) (attachments []*schemas.FileAttachment) {
	for _, event := range events {
		for _, eventAttachment := range event.Attachments {
			// skip pdf files
			if eventAttachment.MimeType == "application/pdf" {
				logrus.Debugf("skipping pdf: %s@%s \n", event.ChildName, event.EventTime)
				continue
			}
			att := schemas.NewFileAttachment(event, eventAttachment)
			attachments = append(attachments, att)
		}
	}
	return attachments
}
