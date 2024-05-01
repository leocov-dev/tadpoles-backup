package bright_horizons

import (
	"net/url"
	"time"
)

type Endpoints struct {
	root        *url.URL
	apiV2Root   *url.URL
	loginUrl    *url.URL
	validateUrl *url.URL
	profileUrl  *url.URL
}

func newEndpoints() Endpoints {
	loginUrl, _ := url.Parse("https://bhlogin.brighthorizons.com")
	//bhApiUrl, _ := url.Parse("https://mbdwgateway.brighthorizons.com/api")
	rootUrl, _ := url.Parse("https://mybrightday.brighthorizons.com")
	apiV2Root := rootUrl.JoinPath("api", "v2")

	return Endpoints{
		root:        rootUrl,
		apiV2Root:   apiV2Root,
		loginUrl:    loginUrl,
		validateUrl: apiV2Root.JoinPath("auth", "jwt", "validate"),
		profileUrl:  apiV2Root.JoinPath("user", "profile"),
	}
}

func (e Endpoints) dependentsUrl(userId string) *url.URL {
	return e.apiV2Root.JoinPath("dependents", "guardian", userId)
}

func (e Endpoints) reportsUrl(childId string, start, end time.Time) *url.URL {
	eventsUrl := e.apiV2Root.JoinPath("dependent", childId, "daily_reports")

	eventsUrl.RawQuery = url.Values{
		"start": {start.Format("2006-01-02")},
		"end":   {end.Format("2006-01-02")},
	}.Encode()

	return eventsUrl
}

func (e Endpoints) getChunkedReportUrls(dependentId string, from, to time.Time) (reportUrls []*url.URL) {
	chunk := 0
	isLastChunk := false

	// 10 days per chunk
	chunkDelta := 10 * 24 * time.Hour

	chunkStart := from
	chunkEnd := chunkStart.Add(chunkDelta)

	for !isLastChunk {
		reportUrls = append(reportUrls, e.reportsUrl(dependentId, chunkStart, chunkEnd))

		chunk += 1
		chunkStart = chunkEnd
		chunkEnd = chunkEnd.Add(chunkDelta)
		if chunkEnd.After(to) {
			isLastChunk = true
		}
	}

	return reportUrls
}

func (e Endpoints) MediaUrl(attachmentId string) *url.URL {
	return e.apiV2Root.JoinPath("media", attachmentId)
}
