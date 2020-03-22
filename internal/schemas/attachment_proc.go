package schemas

import (
	"context"
	"fmt"
)

type AttachmentProc struct {
	BackupRoot     string
	Context        context.Context
	ErrorChannel   chan string
	FileAttachment *FileAttachment
}

func NewAttachmentProc(attachment *FileAttachment, backupRoot string, errorChannel chan string, ctx context.Context) *AttachmentProc {
	attClone := attachment
	proc := &AttachmentProc{
		BackupRoot:     backupRoot,
		Context:        ctx,
		ErrorChannel:   errorChannel,
		FileAttachment: attClone,
	}

	return proc
}

func (proc *AttachmentProc) ExecSave() {
	saveName := proc.FileAttachment.GetSaveName()

	select {
	case <-proc.Context.Done():
		proc.ErrorChannel <- fmt.Sprintf("Save canceled: %s", saveName)
		return
	default:
		err := proc.FileAttachment.Save(proc.BackupRoot)
		if err != nil {
			proc.ErrorChannel <- fmt.Sprintf("Failed to save attachment -> %s, Msg: %s", saveName, err.Error())
			return
		}
	}
}

func (proc *AttachmentProc) ExecDownload() {
	saveName := proc.FileAttachment.GetSaveName()

	select {
	case <-proc.Context.Done():
		proc.ErrorChannel <- fmt.Sprintf("Donwload canceled: %s", saveName)
		return
	default:
		err := proc.FileAttachment.Download()
		if err != nil {
			proc.ErrorChannel <- fmt.Sprintf("Failed to download attachment -> %s, %s", saveName, err.Error())
			return
		}
	}
}
