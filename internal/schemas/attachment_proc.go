package schemas

import (
	"context"
	"fmt"
	"tadpoles-backup/internal/utils/progress"
)

type AttachmentProc struct {
	backupRoot     string
	ctx            context.Context
	errorChannel   chan string
	fileAttachment *FileAttachment
	barWrapper     *progress.BarWrapper
}

func NewAttachmentProc(
	attachment *FileAttachment,
	backupRoot string,
	errorChannel chan string,
	ctx context.Context,
	barWrapper *progress.BarWrapper,
) *AttachmentProc {
	proc := &AttachmentProc{
		backupRoot:     backupRoot,
		ctx:            ctx,
		errorChannel:   errorChannel,
		fileAttachment: attachment,
		barWrapper:     barWrapper,
	}

	return proc
}

func (proc *AttachmentProc) Execute() {
	saveName := proc.fileAttachment.SaveName()

	defer func() {
		if proc.barWrapper != nil {
			proc.barWrapper.Increment()
		}
	}()

	select {
	case <-proc.ctx.Done():
		proc.errorChannel <- fmt.Sprintf("Donwload canceled: %s", saveName)
		return
	default:
		err := proc.fileAttachment.Download()
		if err != nil {
			proc.errorChannel <- fmt.Sprintf("Failed to download attachment -> %s, %s", saveName, err.Error())
			return
		}
	}

	select {
	case <-proc.ctx.Done():
		proc.errorChannel <- fmt.Sprintf("Save canceled: %s", saveName)
		return
	default:
		err := proc.fileAttachment.Save(proc.backupRoot)
		if err != nil {
			proc.errorChannel <- fmt.Sprintf("Failed to save attachment -> %s, Msg: %s", saveName, err.Error())
			return
		}
	}
}
