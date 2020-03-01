package api

import (
	"github.com/leocov-dev/tadpoles-backup/internal/client"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Attachment(eventKey string, attachmentKey string) (data []byte, err error) {
	params := url.Values{
		"obj": {eventKey},
		"key": {attachmentKey},
	}

	urlBase, _ := url.Parse(client.AttachmentsEndpoint)
	urlBase.RawQuery = params.Encode()

	resp, err := client.ApiClient.Get(urlBase.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, client.NewRequestError(resp)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return body, nil
}
