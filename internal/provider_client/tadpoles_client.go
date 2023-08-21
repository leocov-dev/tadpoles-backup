package provider_client

import (
	"context"
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

func (c *TadpolesClient) GetAllMediaFiles(ctx context.Context, start, end time.Time, useCache bool) (mediaFiles schemas.MediaFiles, err error) {
	var events tadpoles.Events

	if useCache {
		cachedEvents, err := c.cache.ReadEventCache()
		if err != nil {
			return nil, err
		}
		events = append(events, cachedEvents...)

		cachedEventsLen := len(cachedEvents)
		if cachedEventsLen > 0 {
			lastCachedEvent := cachedEvents[cachedEventsLen-1]
			lastEventTime := lastCachedEvent.EventTime.Time()
			// TODO this falls apart if `end` is not after `lastEventTime`
			start = lastEventTime.Add(1 * time.Second)
		}
	}

	newEvents, err := c.spec.GetEvents(ctx, start, end)
	if err != nil {
		return nil, err
	}
	if len(newEvents) > 0 {
		if useCache {
			err = c.cache.UpdateEventCache(newEvents)
			if err != nil {
				return nil, err
			}
		}
		events = append(events, newEvents...)
	}

	for _, event := range events {
		eventFiles, err := c.spec.GetEventMediaFiles(*event)
		if err != nil {
			return nil, err
		}
		mediaFiles = append(mediaFiles, eventFiles...)
	}

	return mediaFiles, nil
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

func (c *TadpolesClient) ShouldUseCache(_ string) bool {
	return true
}
