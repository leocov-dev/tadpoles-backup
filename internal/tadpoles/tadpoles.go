package tadpoles

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"tadpoles-backup/internal/api"
	"tadpoles-backup/internal/cache"
	"tadpoles-backup/internal/schemas"
	"tadpoles-backup/internal/utils"
	"tadpoles-backup/internal/utils/progress"
	"time"
)

func GetAllEvents() (events api.Events, err error) {
	// start at "zero" time
	lastCachedTime := time.Time{}

	cachedEvents, err := cache.ReadEventCache()
	if err != nil {
		return nil, err
	}

	if len(cachedEvents) > 0 {
		lastCachedTime = cachedEvents[len(cachedEvents)-1].EventTime.Time()
		log.Debugf("lastCachedTime: %s\n", lastCachedTime)

		// add one second to prevent fetch of existing event
		lastCachedTime = lastCachedTime.Add(1 * time.Second)
		events = append(events, cachedEvents...)
	}

	newEvents, err := api.Spec.GetEvents(lastCachedTime, time.Now())
	if err != nil {
		return nil, err
	}

	if len(newEvents) > 0 {
		err = cache.UpdateEventCache(newEvents)
		if err != nil {
			return nil, err
		}

		events = append(events, newEvents...)
	}

	return events, nil
}

func GetEventFileAttachmentData(events api.Events) (fileAttachments schemas.FileAttachments, err error) {
	fileAttachments = eventsToFileAttachments(events)

	return fileAttachments, nil
}

func DownloadFileAttachments(
	newAttachments schemas.FileAttachments,
	backupRoot string,
	ctx context.Context,
	concurrencyLimit int,
	barWrapper *progress.BarWrapper,
) []string {

	errorChan := make(chan string)

	downloadPool := schemas.NewDownloadPool(concurrencyLimit)

	for _, attachment := range newAttachments {
		proc := schemas.NewAttachmentProc(attachment, backupRoot, errorChan, ctx, barWrapper)
		downloadPool.Add(proc)
	}
	downloadPool.Process()

	close(errorChan)

	var saveErrors []string
	for s := range errorChan {
		saveErrors = append(saveErrors, s)
	}

	return saveErrors
}

func GroupAttachmentsByType(attachments schemas.FileAttachments) schemas.FileAttachmentMap {
	attachmentTypeMap := make(map[string]schemas.FileAttachments)

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

func PruneAlreadyDownloaded(
	attachments schemas.FileAttachments,
	backupTarget string,
) (newAttachments schemas.FileAttachments, err error) {
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

func eventsToFileAttachments(events []*api.Event) (attachments schemas.FileAttachments) {
	for _, event := range events {
		for _, eventAttachment := range event.Attachments {
			// skip pdf files
			if eventAttachment.MimeType == "application/pdf" {
				log.Debugf("skipping pdf: %s @ %s \n", event.ChildName, event.EventTime)
				continue
			}
			att := schemas.NewFileAttachment(event, eventAttachment)
			attachments = append(attachments, att)
		}
	}
	return attachments
}

func PrintErrorList(errorMsgs []string) {
	if errorMsgs != nil {
		utils.WriteError("Errors", "")
		for i, e := range errorMsgs {
			utils.WriteErrorSub.Write(fmt.Sprint(i+1), e)
		}
		fmt.Println("")
	}
}
