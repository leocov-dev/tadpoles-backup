package api

import "net/url"

var (
	tadpolesUrl, _ = url.Parse("https://www.tadpoles.com")
	apiV1Root      = tadpolesUrl.JoinPath("remote", "v1")
)

type Endpoints struct {
	Events      *url.URL
	Attachments *url.URL
	Parameters  *url.URL
}

func newEndpoints() Endpoints {
	return Endpoints{
		Events:      apiV1Root.JoinPath("events"),
		Attachments: apiV1Root.JoinPath("obj_attachment"),
		Parameters:  apiV1Root.JoinPath("parameters"),
	}
}
