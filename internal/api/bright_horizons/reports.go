package bright_horizons

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"sort"
	"tadpoles-backup/internal/utils"
	"time"
)

type Report struct {
	Id        string `json:"id"`
	Dependent *Child
	Created   time.Time   `json:"created"`
	Snapshots []*Snapshot `json:"snapshot_entries"`
}
type Reports []*Report

func (r Reports) Dedupe() Reports {
	filtered := Reports{}
	foundMap := make(map[string]bool)

	for _, report := range r {
		if _, found := foundMap[report.Id]; found {
			log.Debug("duplicate report ", report.Id)
		}
		foundMap[report.Id] = true
		filtered = append(filtered, report)
	}

	return filtered
}

type By func(e1, e2 *Report) bool

func (r Reports) Sort(by By) {
	es := &reportSorter{
		reports: r,
		by:      by,
	}

	sort.Sort(es)
}

type reportSorter struct {
	reports Reports
	by      By
}

func (s *reportSorter) Len() int {
	return len(s.reports)
}

func (s *reportSorter) Swap(i, j int) {
	s.reports[i], s.reports[j] = s.reports[j], s.reports[i]
}

func (s *reportSorter) Less(i, j int) bool {
	return s.by(s.reports[i], s.reports[j])
}

func ByReportDate(r1, r2 *Report) bool {
	return r1.Created.Before(r2.Created)
}

func fetchDependentReports(client *http.Client, reportUrl *url.URL, dependent Child) (reports Reports, err error) {
	resp, err := client.Get(reportUrl.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, utils.NewRequestError(resp, "Failed to fetch reports")
	}

	defer utils.CloseWithLog(resp.Body)
	body, _ := io.ReadAll(resp.Body)

	err = json.Unmarshal(body, &reports)

	for _, r := range reports {
		r.Dependent = &dependent
	}

	return reports, err
}
