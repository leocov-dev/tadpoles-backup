package provider_client

import (
	"context"
	"net/http"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/schemas"
	"time"
)

type ProviderClient interface {
	LoginIfNeeded() error
	GetAccountInfo() (*schemas.AccountInfo, error)
	GetAllMediaFiles(ctx context.Context, start, end time.Time, useCache bool) (schemas.MediaFiles, error)
	ClearLoginData() error
	ClearCache() error
	ClearAll() []error
	GetHttpClient() *http.Client
	ShouldUseCache(operation string) bool
	ResetUserPassword(email string) error
}

func GetProviderClient() ProviderClient {
	switch config.Provider.String() {
	case config.BrightHorizons:
		return NewBrightHorizonsClient()
	default:
		return NewTadpolesClient()
	}
}
