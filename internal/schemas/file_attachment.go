package schemas

import (
	"errors"
	"fmt"
	"github.com/h2non/filetype"
	"github.com/leocov-dev/tadpoles-backup/internal/api"
	"github.com/leocov-dev/tadpoles-backup/internal/utils"
	log "github.com/sirupsen/logrus"
	"io"
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
	tempFileName  string
	backupTarget  string
	Exists        bool
}

// set the parent directory for the save target
func (a *FileAttachment) SetBackupRoot(backupTarget string) {
	a.backupTarget = backupTarget
}

// get the save target file name without extension
func (a *FileAttachment) GetSaveName() string {
	timestamp := fmt.Sprintf("%d%02d%02d%02d%02d%02d",
		a.EventTime.Year(),
		a.EventTime.Month(),
		a.EventTime.Day(),
		a.EventTime.Hour(),
		a.EventTime.Minute(),
		a.EventTime.Second(),
	)
	return fmt.Sprintf("%s_%s", timestamp, a.ChildName)
}

// get the save target directory
func (a *FileAttachment) GetSaveDir() string {
	return path.Join(a.backupTarget, fmt.Sprint(a.EventTime.Year()), fmt.Sprintf("%d-%02d-%02d", a.EventTime.Year(), a.EventTime.Month(), a.EventTime.Day()))
}

// get the path and filename for the final save location
func (a *FileAttachment) GetSaveTarget() (filePath string, err error) {
	if a.backupTarget == "" {
		return "", errors.New("backup target must be set before writing")
	}

	kind, err := filetype.MatchFile(a.tempFileName)
	if err != nil {
		return "", err
	}

	dir := a.GetSaveDir()
	fileName := fmt.Sprintf("%s.%s", a.GetSaveName(), kind.Extension)
	return path.Join(dir, fileName), nil
}

// download file to a temporary directory
func (a *FileAttachment) Download() (err error) {
	log.Debug("Downloading: ", a.AttachmentKey)

	resp, err := api.Attachment(a.EventKey, a.AttachmentKey)
	if err != nil {
		return err
	}

	tempFile, err := ioutil.TempFile("", "tpbk_*")
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return nil
	}

	a.tempFileName = tempFile.Name()
	return nil
}

// create the necessary directories and move the temporary file to the target with a new name
func (a *FileAttachment) Save() (err error) {

	savePath, err := a.GetSaveTarget()
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Dir(savePath), os.ModePerm)
	if err != nil {
		return err
	}

	log.Debug("Saving to: ", savePath)
	err = utils.MoveFile(a.tempFileName, savePath)
	if err != nil {
		return err
	}

	return nil
}
