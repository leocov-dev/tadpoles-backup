package schemas

import (
	"errors"
	"fmt"
	"github.com/h2non/filetype"
	"github.com/leocov-dev/tadpoles-backup/internal/api"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type FileAttachment struct {
	Comment       string
	AttachmentKey string
	EventKey      string
	ChildName     string
	CreateTime    time.Time
	EventTime     time.Time
	bytes         []byte
	backupTarget  string
	Exists        bool
}

func (a *FileAttachment) SetBackupDir(backupTarget string) {
	a.backupTarget = backupTarget
}

func (a *FileAttachment) GetSaveName() string {
	timestamp := fmt.Sprintf("%d%d%d%d%d%d",
		a.EventTime.Year(),
		a.EventTime.Month(),
		a.EventTime.Day(),
		a.EventTime.Hour(),
		a.EventTime.Minute(),
		a.EventTime.Second(),
	)
	return fmt.Sprintf("%s_%s", timestamp, a.ChildName)
}

func (a *FileAttachment) GetSaveDir() string {
	return path.Join(a.backupTarget, fmt.Sprintf("%d-%02d", a.EventTime.Year(), a.EventTime.Month()))
}

func (a *FileAttachment) GetSavePath() (filePath string, err error) {
	if a.backupTarget == "" {
		return "", errors.New("backup target must be set before writing")
	}

	kind, err := filetype.Match(a.bytes)
	if err != nil {
		return "", err
	}

	dir := a.GetSaveDir()
	fileName := fmt.Sprintf("%s.%s", a.GetSaveName(), kind.Extension)
	return path.Join(dir, fileName), nil
}

func (a *FileAttachment) Download() (err error) {
	log.Debug("Downloading: ", a.AttachmentKey)
	data, err := api.Attachment(a.EventKey, a.AttachmentKey)
	if err != nil {
		return err
	}
	a.bytes = data
	return nil
}

func (a *FileAttachment) Write() (err error) {

	savePath, err := a.GetSavePath()
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Dir(savePath), os.ModePerm)
	if err != nil {
		return err
	}

	log.Debug("Saving to: ", savePath)
	err = ioutil.WriteFile(savePath, a.bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
