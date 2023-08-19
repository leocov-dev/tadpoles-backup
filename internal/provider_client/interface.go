package provider_client

import (
	"net/http"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/schemas"
	"time"
)

type ProviderClient interface {
	LoginIfNeeded() error
	GetAccountInfo() (*schemas.AccountInfo, error)
	GetAllMediaFiles(start, end time.Time) (schemas.MediaFiles, error)
	ClearLoginData() error
	ClearCache() error
	ClearAll() []error
	GetHttpClient() *http.Client
}

func GetProviderClient() ProviderClient {
	switch config.Provider.String() {
	case config.BRIGHT_HORIZONS:
		return NewBrightHorizonsClient()
	default:
		return NewTadpolesClient()
	}
}
