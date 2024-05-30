package bright_horizons

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/api/tadpoles"
	"tadpoles-backup/internal/http_utils"
	"tadpoles-backup/internal/schemas"
	"time"
)

type Provider struct {
	httpClient *http.Client
	e          endpoints
	c          *tadpoles.ApiCache
}

func NewProvider() *Provider {
	ep := newEndpoints()
	cache := tadpoles.NewApiCache(config.BrightHorizons)
	return &Provider{
		httpClient: &http.Client{
			Jar:       http_utils.DeserializeCookies(cache.CookieFile, ep.root),
			Transport: &http_utils.RandomUserAgentTransport{},
			Timeout:   60 * time.Second,
		},
		e: ep,
		c: cache,
	}
}

func (p *Provider) HttpClient() *http.Client {
	return p.httpClient
}

func (p *Provider) LoginIfNeeded() error {
	// TODO
	return nil
}

func (p *Provider) FetchAccountInfo() (*schemas.AccountInfo, error) {
	// TODO
	return nil, nil
}

func (p *Provider) FetchAllMediaFiles(ctx context.Context, start, end time.Time) (schemas.MediaFiles, error) {
	// TODO
	return nil, nil
}

func (p *Provider) ClearLoginData() error {
	return p.c.ClearLoginCookie()
}

func (p *Provider) ClearCache() error {
	return p.c.ClearCache()
}

func (p *Provider) ClearAll() []error {
	cookieErr := p.c.ClearLoginCookie()
	cacheErr := p.c.ClearCache()

	return []error{cookieErr, cacheErr}
}

func (p *Provider) ResetUserPassword(email string) error {
	return errors.New(fmt.Sprintf("%s client does not support password reset", config.Provider.String()))
}
