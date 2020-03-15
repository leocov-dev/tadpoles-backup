package tadpoles

import (
	"github.com/gosuri/uiprogress"
	"github.com/leocov-dev/tadpoles-backup/internal/api"
	"github.com/leocov-dev/tadpoles-backup/internal/schemas"
	"github.com/leocov-dev/tadpoles-backup/internal/utils"
	"sync"
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

func DownloadFileAttachments(attachments []*schemas.FileAttachment, backupTarget string) ([]string, error) {
	err := checkAlreadyDownloaded(attachments, backupTarget)
	if err != nil {
		return nil, err
	}

	uiprogress.Start()
	dwnd := uiprogress.AddBar(len(attachments)).
		AppendCompleted().
		PrependFunc(func(b *uiprogress.Bar) string {
			return utils.HiCyan.Sprint("Downloading")
		})

	wg := &sync.WaitGroup{}
	errorChan := make(chan string, len(attachments))

	for _, attachment := range attachments {
		if attachment.AlreadyDownloaded {
			dwnd.Incr()
			//log.Debug("Already exists: ", attachment.GetSaveName())
			continue
		}

		wg.Add(1)
		attachment.SetBackupRoot(backupTarget)

		go saveFileAttachment(attachment, wg, dwnd, errorChan)
	}
	wg.Wait()
	close(errorChan)

	var saveErrors []string
	for s := range errorChan {
		saveErrors = append(saveErrors, s)
	}

	return saveErrors, nil
}
