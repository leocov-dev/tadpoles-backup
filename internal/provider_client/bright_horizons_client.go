package provider_client

import (
	"fmt"
	"net/http"
	"tadpoles-backup/internal/api/bright_horizons"
	"tadpoles-backup/internal/cache"
	"tadpoles-backup/internal/schemas"
	"tadpoles-backup/internal/user_input"
	"time"
)

type BrightHorizonsClient struct {
	spec  *bright_horizons.ApiSpec
	cache *cache.BrightHorizonsCache
}

func NewBrightHorizonsClient() *BrightHorizonsClient {
	return &BrightHorizonsClient{
		spec:  bright_horizons.NewApiSpec(),
		cache: cache.NewBrightHorizonsCache(),
	}
}

func (c *BrightHorizonsClient) GetHttpClient() *http.Client {
	return c.spec.Client
}

func (c *BrightHorizonsClient) LoginIfNeeded() error {
	if c.spec.NeedsLogin(c.cache.GetApiKey()) {
		username, password := user_input.GetUsernameAndPassword()

		apiKey, err := c.spec.DoLogin(username, password)
		if err != nil {
			return err
		}

		err = c.cache.StoreApiKey(apiKey)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *BrightHorizonsClient) GetAccountInfo() (info *schemas.AccountInfo, err error) {
	dependents, err := c.spec.GetAccountData()
	if err != nil {
		return nil, err
	}

	info = &schemas.AccountInfo{
		FirstEvent: time.Now(), // initialize to a later time than expected
	}

	for _, dep := range dependents {
		info.Dependants = append(
			info.Dependants,
			fmt.Sprintf("%s %s", dep.FirstName, dep.LastName),
		)

		if dep.FirstRecord.Before(info.FirstEvent) {
			info.FirstEvent = dep.FirstRecord
		}
		if dep.LastRecord.After(info.LastEvent) {
			info.LastEvent = dep.LastRecord
		}
	}

	return info, nil
}

func (c *BrightHorizonsClient) GetAllMediaFiles(start, end time.Time) (attachments schemas.MediaFiles, err error) {
	return nil, nil
}

func (c *BrightHorizonsClient) ClearLoginData() error {
	return c.cache.ClearLoginData()
}

func (c *BrightHorizonsClient) ClearCache() error {
	return c.cache.ClearCache()
}

func (c *BrightHorizonsClient) ClearAll() []error {
	cookieErr := c.cache.ClearLoginData()
	cacheErr := c.cache.ClearCache()

	return []error{cookieErr, cacheErr}
}
