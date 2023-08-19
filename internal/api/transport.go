package api

import (
	"github.com/corpix/uarand"
	"net/http"
)

type RandomUserAgentTransport struct{}

func (t *RandomUserAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", uarand.GetRandom())
	return http.DefaultTransport.RoundTrip(req)
}
