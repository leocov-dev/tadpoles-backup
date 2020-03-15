package tadpoles

import (
	"fmt"
	"github.com/gosuri/uiprogress"
	"github.com/leocov-dev/tadpoles-backup/internal/api"
	"github.com/leocov-dev/tadpoles-backup/internal/schemas"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"strings"
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

// update FileAttachment.AlreadyDownloaded based on the existence of a matching file name, minus extension
func checkAlreadyDownloaded(attachments []*schemas.FileAttachment, backupTarget string) (err error) {
	attachmentNames := make(map[string]int)
	for i, att := range attachments {
		attachmentNames[att.GetSaveName()] = i
	}

	err = filepath.Walk(backupTarget,
		func(path_ string, info_ os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info_.IsDir() {
				return nil
			}

			file := path.Base(path_)
			minusExtension := strings.TrimSuffix(file, filepath.Ext(file))

			if i, ok := attachmentNames[minusExtension]; ok {
				attachments[i].AlreadyDownloaded = true
			}

			return nil
		})

	return err
}

// goroutine to download the file attachment and save it
func saveFileAttachment(attachment *schemas.FileAttachment, progress *uiprogress.Bar, c chan string) {
	err := attachment.Download()
	if err != nil {
		c <- fmt.Sprintf("Failed to download attachment -> %s, %s", attachment.GetSaveName(), err.Error())
		return
	}
	err = attachment.Save()
	if err != nil {
		c <- fmt.Sprintf("Failed to save attachment -> %s, Msg: %s", attachment.GetSaveName(), err.Error())
		return
	}
	if progress != nil {
		progress.Incr()
	}
}
