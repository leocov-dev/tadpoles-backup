package api

import (
	"github.com/leocov-dev/tadpoles-backup/internal/client"
	"net/http"
	"net/url"
)

func GetAttachment(eventKey string, attachmentKey string) (resp *http.Response, err error) {
	params := url.Values{
		"obj": {eventKey},
		"key": {attachmentKey},
	}

	urlBase, _ := url.Parse(client.AttachmentsEndpoint)
	urlBase.RawQuery = params.Encode()

	resp, err = client.ApiClient.Get(urlBase.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, client.NewRequestError(resp)
	}

	return resp, err
}
