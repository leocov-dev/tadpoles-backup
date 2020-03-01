package tadpoles_api

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
	"sync"
	"time"
)

var fileCache map[string]bool

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

	return eventsToAttachments(events), nil
}

func DownloadFileAttachments(attachments []*schemas.FileAttachment, backupTarget string) error {
	fmt.Println("Downloads Started...")

	err := tagExisting(attachments, backupTarget)
	if err != nil {
		return err
	}

	uiprogress.Start()
	dwnd := uiprogress.AddBar(len(attachments))
	dwnd.AppendCompleted()

	var wg = &sync.WaitGroup{}

	for _, attachment := range attachments {
		if attachment.Exists {
			dwnd.Incr()
			log.Debug("Already exists: ", attachment.GetSaveName())
			continue
		}

		wg.Add(1)
		attachment.SetBackupDir(backupTarget)

		go saveFileAttachment(attachment, wg, dwnd)
	}
	wg.Wait()
	return nil
}

func tagExisting(attachments []*schemas.FileAttachment, backupTarget string) (err error) {

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

			if val, ok := attachmentNames[minusExtension]; ok {
				attachments[val].Exists = true
			}

			return nil
		})

	return err
}

func saveFileAttachment(attachment *schemas.FileAttachment, group *sync.WaitGroup, progress *uiprogress.Bar) {
	defer group.Done()

	err := attachment.Download()
	if err != nil {
		log.Errorf("Failed to download attachment")
		return
	}
	err = attachment.Write()
	if err != nil {
		log.Errorf("Failed to save attachment")
		return
	}
	progress.Incr()
}
