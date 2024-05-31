package schemas

import (
	"context"
	"tadpoles-backup/internal/http_utils"
	"tadpoles-backup/internal/interfaces"
	"time"
)

type Provider interface {
	HttpClient() interfaces.HttpClient
	LoginIfNeeded() error
	FetchAccountInfo() (*AccountInfo, error)
	FetchAllMediaFiles(ctx context.Context, start, end time.Time) (http_utils.MediaFiles, error)
	ClearLoginData() error
	ClearCache() error
	ClearAll() []error
	ResetUserPassword(email string) error
}
