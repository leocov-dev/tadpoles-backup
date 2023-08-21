package provider_client

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"tadpoles-backup/internal/api/bright_horizons"
	"tadpoles-backup/internal/async"
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
		info.Dependants = append(info.Dependants, dep.DisplayName())

		if dep.FirstRecord.Before(info.FirstEvent) {
			info.FirstEvent = dep.FirstRecord
		}
		if dep.LastRecord.After(info.LastEvent) {
			info.LastEvent = dep.LastRecord
		}
	}

	return info, nil
}

type reportTaskResult struct {
	dependentId string
	reports     bright_horizons.Reports
}

type reportTask struct {
	apiSpec     *bright_horizons.ApiSpec
	resultsChan chan reportTaskResult
	dependent   bright_horizons.Dependent
	reportUrl   *url.URL
}

func (t *reportTask) Run() error {
	log.Debug("fetching chunk...")
	reports, err := t.apiSpec.GetReportsChunk(t.dependent, t.reportUrl)
	if err != nil {
		return err
	}

	filteredReports := bright_horizons.Reports{}

	for _, r := range reports {
		if len(r.Snapshots) > 0 {
			for _, s := range r.Snapshots {
				err = s.HydrateMediaData(t.apiSpec.Client, t.apiSpec.Endpoints.MediaUrl)
				if err != nil {
					return err
				}
			}
			filteredReports = append(filteredReports, r)
		}
	}

	if len(filteredReports) > 0 {
		log.Debug("sending chunk result ", len(filteredReports))
		t.resultsChan <- reportTaskResult{
			dependentId: t.dependent.Id,
			reports:     filteredReports,
		}
	}

	log.Debug("fetching chunk Done")
	return nil
}

func (c *BrightHorizonsClient) GetAllMediaFiles(ctx context.Context, start, end time.Time, useCache bool) (mediaFiles schemas.MediaFiles, err error) {
	dependents, err := c.spec.GetAccountData()
	if err != nil {
		return nil, err
	}

	var combinedReports bright_horizons.Reports
	depUrlMap := make(map[string][]*url.URL)
	resultCount := 0

	for _, d := range dependents {
		if useCache {
			cachedReports, err := c.cache.ReadReportCache(d.Id)
			if err != nil {
				return nil, err
			}
			combinedReports = append(combinedReports, cachedReports...)

			cacheLen := len(cachedReports)
			if cacheLen > 0 {
				lastCachedReport := cachedReports[cacheLen-1]
				lastReportTime := lastCachedReport.Created
				// TODO this falls apart if `end` is not after `lastReportTime`
				start = lastReportTime.Add(24 * time.Hour)
			}
		}

		log.Debugf("Schedule report fetch for: %s", d.DisplayName())
		reportUrls := c.spec.GetReportUrls(d, start, end)
		resultCount += len(reportUrls)
		depUrlMap[d.Id] = reportUrls
	}

	// MUST be a buffered channel since we don't pull results until
	// the end!
	resultsChan := make(chan reportTaskResult, resultCount)

	taskPool := async.NewTaskPool(
		ctx,
		func() { close(resultsChan) },
	)

	for _, d := range dependents {
		for i, chunk := range depUrlMap[d.Id] {
			log.Debugf("Schedule Chunk: %d - %s", i, chunk.Query())
			err = taskPool.AddTask(&reportTask{
				apiSpec:     c.spec,
				resultsChan: resultsChan,
				dependent:   d,
				reportUrl:   chunk,
			})
			if err != nil {
				return nil, err
			}
		}
	}
	taskPool.Wait()

	taskErrors := taskPool.Errors()
	if taskErrors != nil {
		return nil, taskErrors
	}

	log.Debug("report task pool Done")

	for chunk := range resultsChan {
		select {
		case <-ctx.Done():
			return nil, errors.New("canceled")
		default:
			if useCache {
				err = c.cache.UpdateReportCache(chunk.dependentId, chunk.reports)
				if err != nil {
					return nil, err
				}
			}
			log.Debug("combine results chunk ", len(chunk.reports))
			combinedReports = append(combinedReports, chunk.reports...)
		}
	}

	combinedReports = combinedReports.Dedupe()
	combinedReports.Sort(bright_horizons.ByReportDate)

	log.Debug("report collection Done ", len(combinedReports))

	for _, r := range combinedReports {
		for _, s := range r.Snapshots {
			mediaFiles = append(mediaFiles, bright_horizons.NewMediaFileFromReportSnapshot(*r, *s))
		}
	}

	log.Debug("media files compiled ", len(mediaFiles))

	return mediaFiles, nil
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

func (c *BrightHorizonsClient) ShouldUseCache(operation string) bool {
	switch operation {
	case "backup":
		return false
	default:
		return true
	}
}
