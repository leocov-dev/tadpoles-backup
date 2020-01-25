package tadpoles_api

import (
	"github.com/leocov-dev/tadpoles-backup/internal/client"
	"github.com/leocov-dev/tadpoles-backup/internal/schemas"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
)

func ApiAttachment(attachment *schemas.FileAttachment) (data []byte, err error) {
	log.Debug("Get attachment", attachment)
	params := url.Values{
		"obj": {attachment.EventKey},
		"key": {attachment.AttachmentKey},
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
