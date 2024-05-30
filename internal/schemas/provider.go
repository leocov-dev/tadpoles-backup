package schemas

import (
	"context"
	"net/http"
	"time"
)

type Provider interface {
	HttpClient() *http.Client
	LoginIfNeeded() error
	FetchAccountInfo() (*AccountInfo, error)
	FetchAllMediaFiles(ctx context.Context, start, end time.Time) (MediaFiles, error)
	ClearLoginData() error
	ClearCache() error
	ClearAll() []error
	ResetUserPassword(email string) error
}
