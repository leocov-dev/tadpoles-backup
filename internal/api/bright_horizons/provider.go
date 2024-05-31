package bright_horizons

import (
	"context"
	"errors"
	"fmt"
	"github.com/corpix/uarand"
	log "github.com/sirupsen/logrus"
	"net/http"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/api/tadpoles"
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
	c          *tadpoles.ApiCache
	// has cached bh auth token
	bhTokenClient *http.Client
	// client without any auth header modifications
	noAuthClient *http.Client
}

func NewProvider() *Provider {
	ep := newEndpoints()
	cache := tadpoles.NewApiCache(config.BrightHorizons)
	cookieJar := http_utils.DeserializeCookies(cache.CookieFile, ep.root)
	timeout := 60 * time.Second
	return &Provider{
		httpClient: &http.Client{
			Jar: cookieJar,
			Transport: &CachedTokenTransport{
				tokenName: config.Tadpoles,
				cache:     cache,
			},
			Timeout: timeout,
		},
		e: ep,
		c: cache,
		bhTokenClient: &http.Client{
			Jar: cookieJar,
			Transport: &CachedTokenTransport{
				tokenName: config.BrightHorizons,
				cache:     cache,
			},
			Timeout: timeout,
		},
		noAuthClient: &http.Client{
			Jar:       cookieJar,
			Transport: &http_utils.RandomUserAgentTransport{},
			Timeout:   timeout,
		},
	}
}

func (p *Provider) HttpClient() interfaces.HttpClient {
	return p.httpClient
}

func (p *Provider) LoginIfNeeded() error {
	// TODO: not sure if we should call a bh endpoint or a tadpoles endpoint...
	checkErr := checkLogin(p.bhTokenClient, p.e.dependentsUrl)

	if checkErr != nil {
		username, password := user_input.GetUsernameAndPassword()

		log.Debug("Login...")
		rvt, rvtErr := fetchRequestVerificationToken(p.noAuthClient, p.e.rvtUrl)
		if rvtErr != nil {
			return rvtErr
		}

		redirect, loginErr := login(p.noAuthClient, p.e.loginUrl, username, password, rvt)
		if loginErr != nil {
			return loginErr
		}

		log.Debug("Login successful, Admit...")
		action, samlResponse, samlErr := startSaml(p.noAuthClient, p.e.samlUrl(redirect))
		if samlErr != nil {
			return samlErr
		}

		bhToken, samlActionErr := finishSaml(p.noAuthClient, action, samlResponse)
		if samlActionErr != nil {
			return samlActionErr
		}
		bhStoreErr := p.c.StoreToken(config.BrightHorizons, bhToken)
		if bhStoreErr != nil {
			return bhStoreErr
		}

		tadpolesToken, tokenErr := exchangeToken(p.noAuthClient, p.e.token2Url, bhToken)
		if tokenErr != nil {
			return tokenErr
		}
		// it's not clear we need to store and reuse the tadpoles token for
		// later API calls when fetching data
		tpStoreErr := p.c.StoreToken(config.Tadpoles, tadpolesToken)
		if tpStoreErr != nil {
			return tpStoreErr
		}

		cookies, jwtAdmitErr := admitRedirect(p.noAuthClient, p.e.jwtRedirectUrl(tadpolesToken))
		if jwtAdmitErr != nil {
			return jwtAdmitErr
		}

		expires, serializeErr := http_utils.SerializeCookies(p.c.CookieFile, cookies)
		if serializeErr != nil {
			return serializeErr
		}

		if config.IsHumanReadable() {
			utils.WriteInfo(
				"Login expires",
				expires.In(time.Local).Format("Mon Jan 2 03:04:05 PM"),
			)
			fmt.Println("")
		}

		// TODO: not sure what to do about various cookies, the noAuthClient
		//  should have all the cookies so far collected?
		p.bhTokenClient.Jar = p.noAuthClient.Jar
		p.httpClient.Jar = p.noAuthClient.Jar
	}

	return nil
}

func (p *Provider) FetchAccountInfo() (*schemas.AccountInfo, error) {
	// TODO: no idea how this should work, should it call only tadpoles API's
	//  should it call bright_horizons.fetchDependents() (if that is even valid)
	parameters, paramErr := tadpoles.FetchParameters(p.httpClient, p.e.parametersUrl)
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
	// TODO: no idea how this should work....
	allEvents, eventsErr := tadpoles.FetchAllEvents(ctx, p.c, p.httpClient, p.e, start, end, true)
	if eventsErr != nil {
		return nil, eventsErr
	}

	mediaFiles := http_utils.MediaFiles{}

	for _, event := range allEvents {
		eventMediaFiles := make(http_utils.MediaFiles, len(event.Attachments))
		for i, attachment := range event.Attachments {
			eventMediaFiles[i] = tadpoles.NewMediaFileFromEventAttachment(*event, *attachment, p.e)
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
	return errors.New(fmt.Sprintf("%s client does not support password reset", config.Provider.String()))
}

type CachedTokenTransport struct {
	tokenName string
	cache     *tadpoles.ApiCache
}

func (t *CachedTokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", uarand.GetRandom())
	token, err := t.cache.GetToken(t.tokenName)
	if err != nil {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	return http.DefaultTransport.RoundTrip(req)
}
