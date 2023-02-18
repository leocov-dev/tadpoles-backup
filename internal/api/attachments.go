package api

import (
	"net/http"
	"net/url"
	"tadpoles-backup/internal/client"
)

func GetAttachment(eventKey string, attachmentKey string) (resp *http.Response, err error) {
	params := url.Values{
		"obj": {eventKey},
		"key": {attachmentKey},
	}

	urlBase, _ := url.Parse(client.AttachmentsEndpoint)
	urlBase.RawQuery = params.Encode()

	resp, err = client.GetApiClient().Get(urlBase.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, client.NewRequestError(resp)
	}

	return resp, err
}
