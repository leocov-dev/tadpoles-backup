package tadpoles_api

import (
	"github.com/gosuri/uiprogress"
	"github.com/leocov-dev/tadpoles-backup/internal/api"
	"github.com/leocov-dev/tadpoles-backup/internal/schemas"
	"github.com/leocov-dev/tadpoles-backup/internal/utils"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"strings"
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
			log.Debug("Already exists: ", attachment.GetSaveName())
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
func saveFileAttachment(attachment *schemas.FileAttachment, group *sync.WaitGroup, progress *uiprogress.Bar, c chan string) {
	defer group.Done()

	err := attachment.Download()
	if err != nil {
		c <- utils.HiRed.Sprintf("Failed to download attachment -> %s, %s", attachment.GetSaveName(), err.Error())
		return
	}
	err = attachment.Save()
	if err != nil {
		c <- utils.HiRed.Sprintf("Failed to save attachment -> %s, %s", attachment.GetSaveName(), err.Error())
		return
	}
	progress.Incr()
}
