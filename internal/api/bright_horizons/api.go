package bright_horizons

import (
	"github.com/corpix/uarand"
	"net/http"
	"tadpoles-backup/internal/api"
	"time"
)

type ApiSpec struct {
	Client    *http.Client
	endpoints endpoints
}

func NewApiSpec() *ApiSpec {
	return &ApiSpec{
		endpoints: newEndpoints(),
		Client: &http.Client{
			Transport: &api.RandomUserAgentTransport{},
			Timeout:   60 * time.Second,
		},
	}
}

func (a *ApiSpec) setApiKey(apiKey string) {
	a.Client = &http.Client{
		Jar:       a.Client.Jar,
		Transport: &ApiKeyTransport{apiKey: apiKey},
		Timeout:   a.Client.Timeout,
	}
}

func (a *ApiSpec) NeedsLogin(apiKey string) bool {
	a.setApiKey(apiKey)
	resp, err := a.Client.Get(a.endpoints.profileUrl.String())

	return err != nil || resp.StatusCode != http.StatusOK
}

func (a *ApiSpec) DoLogin(username, password string) (string, error) {
	return login(a.Client, a.endpoints.loginUrl, username, password)
}

func (a *ApiSpec) GetAccountData() (dependents DependentResponse, err error) {
	profile, err := fetchProfile(a.Client, a.endpoints.profileUrl)
	if err != nil {
		return nil, err
	}

	dependents, err = fetchDependents(a.Client, a.endpoints.dependentsUrl(profile.UserId))
	if err != nil {
		return nil, err
	}

	return dependents, nil
}

type ApiKeyTransport struct {
	apiKey string
}

func (t *ApiKeyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", uarand.GetRandom())
	req.Header.Add("X-Api-Key", t.apiKey)
	return http.DefaultTransport.RoundTrip(req)
}
