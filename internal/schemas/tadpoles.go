package schemas

import (
	"time"
)

type FileAttachment struct {
	Comment       string
	AttachmentKey string
	EventKey      string
	ChildName     string
	CreateTime    time.Time
	EventTime     time.Time
}

type Info struct {
	FirstEvent time.Time
	LastEvent  time.Time
	Dependants []string
}
