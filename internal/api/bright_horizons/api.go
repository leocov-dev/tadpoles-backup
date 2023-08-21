package bright_horizons

import (
	"errors"
	"github.com/corpix/uarand"
	"net/http"
	"net/url"
	"tadpoles-backup/internal/api"
	"tadpoles-backup/internal/schemas"
	"time"
)

type ApiSpec struct {
	Client    *http.Client
	Endpoints Endpoints
}

func NewApiSpec() *ApiSpec {
	return &ApiSpec{
		Endpoints: newEndpoints(),
		Client: &http.Client{
			Transport: &api.RandomUserAgentTransport{},
			Timeout:   120 * time.Second,
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
	resp, err := a.Client.Get(a.Endpoints.profileUrl.String())

	return err != nil || resp.StatusCode != http.StatusOK
}

func (a *ApiSpec) DoLogin(username, password string) (string, error) {
	return login(a.Client, a.Endpoints.loginUrl, username, password)
}

func (a *ApiSpec) GetAccountData() (dependents Dependents, err error) {
	profile, err := fetchProfile(a.Client, a.Endpoints.profileUrl)
	if err != nil {
		return nil, err
	}

	dependents, err = fetchDependents(a.Client, a.Endpoints.dependentsUrl(profile.UserId))
	if err != nil {
		return nil, err
	}

	return dependents, nil
}

func (a *ApiSpec) GetReportUrls(dependent Dependent, start, end time.Time) (reportUrls []*url.URL) {
	return a.Endpoints.getChunkedReportUrls(dependent.Id, start, end)
}

func (a *ApiSpec) GetReportsChunk(dependent Dependent, reportUrl *url.URL) (reports Reports, err error) {
	reports, err = fetchDependentReports(a.Client, reportUrl, dependent)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

func (a *ApiSpec) GetReportMediaFiles(report Report) (schemas.MediaFiles, error) {
	mediaFiles := make(schemas.MediaFiles, len(report.Snapshots))

	for i, snapshot := range report.Snapshots {
		if snapshot.MediaResponse.MimeType == "" {
			return nil, errors.New("snapshot does not have hydrated media data")
		}
		mediaFiles[i] = NewMediaFileFromReportSnapshot(report, *snapshot)
	}

	return mediaFiles, nil
}

type ApiKeyTransport struct {
	apiKey string
}

func (t *ApiKeyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", uarand.GetRandom())
	req.Header.Add("X-Api-Key", t.apiKey)
	return http.DefaultTransport.RoundTrip(req)
}
