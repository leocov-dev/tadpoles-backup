package schemas

import (
	"fmt"
	"net/http"
	"tadpoles-backup/internal/utils/progress"
)

type DownloadTask struct {
	downloadRoot string
	mediaFile    MediaFile
	client       *http.Client
	progress     *progress.BarWrapper
}

func NewDownloadTask(
	client *http.Client,
	mediaFile MediaFile,
	downloadRoot string,
	sharedProgressBar *progress.BarWrapper,
) *DownloadTask {
	return &DownloadTask{
		downloadRoot: downloadRoot,
		mediaFile:    mediaFile,
		client:       client,
		progress:     sharedProgressBar,
	}
}

func (task *DownloadTask) Run() error {
	defer func() {
		if task.progress != nil {
			task.progress.Increment()
		}
	}()

	err := task.mediaFile.Download(task.client, task.downloadRoot)
	if err != nil {
		return fmt.Errorf(
			"failed to download attachment -> %s, %s",
			task.mediaFile.FileName,
			err.Error(),
		)
	}

	return nil
}
