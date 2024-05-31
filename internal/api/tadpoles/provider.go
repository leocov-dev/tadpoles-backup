package tadpoles

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/http_utils"
	"tadpoles-backup/internal/interfaces"
	"tadpoles-backup/internal/schemas"
	"tadpoles-backup/internal/user_input"
	"tadpoles-backup/internal/utils"
	"time"
)

type Provider struct {
	httpClient *http.Client
	e          endpoints
	c          *ApiCache
}

func NewProvider() *Provider {
	ep := newEndpoints()
	cache := NewApiCache(config.Tadpoles)
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

func (p *Provider) HttpClient() interfaces.HttpClient {
	return p.httpClient
}

func (p *Provider) LoginIfNeeded() error {
	_, checkErr := loginAdmit(p.httpClient, p.e.admitUrl, p.c.CookieFile)

	if checkErr != nil {
		username, password := user_input.GetUsernameAndPassword()

		log.Debug("Login...")
		loginErr := login(p.httpClient, p.e.loginUrl, username, password)
		if loginErr != nil {
			return loginErr
		}

		log.Debug("Login successful, Admit...")
		expires, admitErr := loginAdmit(p.httpClient, p.e.admitUrl, p.c.CookieFile)
		if admitErr != nil {
			return admitErr
		}

		if config.IsHumanReadable() {
			utils.WriteInfo(
				"Login expires",
				expires.In(time.Local).Format("Mon Jan 2 03:04:05 PM"),
			)
			fmt.Println("")
		}
	}

	return nil
}

func (p *Provider) FetchAccountInfo() (*schemas.AccountInfo, error) {
	parameters, paramErr := FetchParameters(p.httpClient, p.e.parametersUrl)
	if paramErr != nil {
		return nil, paramErr
	}

	info := &schemas.AccountInfo{
		FirstEvent: parameters.FirstEventTime.Time(),
		LastEvent:  parameters.LastEventTime.Time(),
	}

	for _, mem := range parameters.Memberships {
		for _, dep := range mem.Dependents {
			info.Dependants = append(info.Dependants, dep.DisplayName)
		}
	}

	return info, nil
}

func (p *Provider) FetchAllMediaFiles(ctx context.Context, start, end time.Time) (http_utils.MediaFiles, error) {
	allEvents, eventsErr := FetchAllEvents(ctx, p.c, p.httpClient, p.e, start, end, true)
	if eventsErr != nil {
		return nil, eventsErr
	}

	mediaFiles := http_utils.MediaFiles{}

	for _, event := range allEvents {
		eventMediaFiles := make(http_utils.MediaFiles, len(event.Attachments))
		for i, attachment := range event.Attachments {
			eventMediaFiles[i] = NewMediaFileFromEventAttachment(*event, *attachment, p.e)
		}
		mediaFiles = append(mediaFiles, eventMediaFiles...)
	}

	return mediaFiles, nil
}

func (p *Provider) ClearLoginData() error {
	return p.c.ClearLoginData()
}

func (p *Provider) ClearCache() error {
	return p.c.ClearCache()
}

func (p *Provider) ClearAll() []error {
	cookieErr := p.c.ClearLoginData()
	cacheErr := p.c.ClearCache()

	return []error{cookieErr, cacheErr}
}

func (p *Provider) ResetUserPassword(email string) error {
	log.Debug("Ask tadpoles.com to reset user: ", email)

	// We must make a new tp_client with the `Host` header set
	// to impersonate www.tadpoles.com or the request will fail.
	// They've implemented a basic security check to limit reset
	// spamming, but luckily they don't validate the header
	// against an IP address etc.
	client := &http.Client{
		Transport: &hostHeaderTransport{
			hostHeader: p.e.root.Host,
		},
		Timeout: 60 * time.Second,
	}

	return requestPasswordReset(client, p.e.resetUrl, email)
}
