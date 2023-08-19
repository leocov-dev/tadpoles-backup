package provider_client

import (
	"fmt"
	"net/http"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/api/tadpoles"
	"tadpoles-backup/internal/cache"
	"tadpoles-backup/internal/schemas"
	"tadpoles-backup/internal/user_input"
	"tadpoles-backup/internal/utils"
	"time"
)

type TadpolesClient struct {
	spec  *tadpoles.ApiSpec
	cache *cache.TadpolesCache
}

func NewTadpolesClient() *TadpolesClient {
	tadpolesCache := cache.NewTadpolesCache()
	return &TadpolesClient{
		spec:  tadpoles.NewApiSpec(tadpolesCache.GetCookieFile()),
		cache: tadpolesCache,
	}
}

func (c *TadpolesClient) GetHttpClient() *http.Client {
	return c.spec.Client
}

func (c *TadpolesClient) LoginIfNeeded() error {
	if c.spec.NeedsLogin(c.cache.GetCookieFile()) {
		username, password := user_input.GetUsernameAndPassword()

		expires, err := c.spec.DoLogin(username, password, c.cache.GetCookieFile())
		if err != nil {
			return err
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

func (c *TadpolesClient) GetAccountInfo() (info *schemas.AccountInfo, err error) {
	parameters, err := c.spec.GetAccountParameters()
	if err != nil {
		return nil, err
	}

	info = &schemas.AccountInfo{
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

func (c *TadpolesClient) GetAllMediaFiles(start, end time.Time) (attachments schemas.MediaFiles, err error) {
	var events tadpoles.Events

	cachedEvents, err := c.cache.ReadEventCache()
	if err != nil {
		return nil, err
	}

	cachedEventsLen := len(cachedEvents)

	if cachedEventsLen > 0 {
		start = cachedEvents[cachedEventsLen-1].EventTime.Time().Add(1 * time.Second)
		events = append(events, cachedEvents...)
	}

	newEvents, err := c.spec.GetEvents(start, end)
	if err != nil {
		return nil, err
	}
	if len(newEvents) > 0 {
		err = c.cache.UpdateEventCache(newEvents)
		if err != nil {
			return nil, err
		}
		events = append(events, newEvents...)
	}

	for _, event := range events {
		attachments = append(attachments, c.spec.GetEventAttachments(event)...)
	}

	return attachments, nil
}

func (c *TadpolesClient) ClearLoginData() error {
	return c.cache.ClearLoginCookie()
}

func (c *TadpolesClient) ClearCache() error {
	return c.cache.ClearCache()
}

func (c *TadpolesClient) ClearAll() []error {
	cookieErr := c.cache.ClearLoginCookie()
	cacheErr := c.cache.ClearCache()

	return []error{cookieErr, cacheErr}
}
