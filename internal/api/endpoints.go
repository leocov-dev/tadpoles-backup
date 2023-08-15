package api

import "net/url"

type Endpoints struct {
	Root        *url.URL
	Events      *url.URL
	Attachments *url.URL
}

func newEndpoints(domain string) Endpoints {
	rootUrl, _ := url.Parse(domain)
	apiV1Root := rootUrl.JoinPath("remote", "v1")
	return Endpoints{
		Root:        rootUrl,
		Events:      apiV1Root.JoinPath("events"),
		Attachments: apiV1Root.JoinPath("obj_attachment"),
	}
}
