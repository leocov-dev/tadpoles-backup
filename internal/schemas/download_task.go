package schemas

import (
	"context"
	"fmt"
	"net/http"
	"tadpoles-backup/internal/utils/progress"
)

type DownloadTask struct {
	downloadRoot string
	ctx          context.Context
	errorChannel chan string
	mediaFile    MediaFile
	client       *http.Client
	progress     *progress.BarWrapper
}

func NewDownloadTask(
	client *http.Client,
	mediaFile MediaFile,
	downloadRoot string,
	errorChannel chan string,
	ctx context.Context,
	sharedProgressBar *progress.BarWrapper,
) *DownloadTask {
	proc := &DownloadTask{
		downloadRoot: downloadRoot,
		ctx:          ctx,
		errorChannel: errorChannel,
		mediaFile:    mediaFile,
		client:       client,
		progress:     sharedProgressBar,
	}

	return proc
}

func (proc *DownloadTask) Run() {
	defer func() {
		if proc.progress != nil {
			proc.progress.Increment()
		}
	}()

	select {
	case <-proc.ctx.Done():
		proc.errorChannel <- fmt.Sprintf(
			"Donwload canceled: %s",
			proc.mediaFile.FileName,
		)
		return
	default:
		err := proc.mediaFile.Download(proc.client, proc.downloadRoot)
		if err != nil {
			proc.errorChannel <- fmt.Sprintf(
				"Failed to download attachment -> %s, %s",
				proc.mediaFile.FileName,
				err.Error(),
			)
			return
		}
	}
}
