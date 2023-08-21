package bright_horizons

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"tadpoles-backup/internal/schemas"
	"tadpoles-backup/internal/utils"
	"time"
)

func NewMediaFileFromReportSnapshot(
	report Report,
	snapshot Snapshot,
) schemas.MediaFile {
	mediaUrl, _ := url.Parse(snapshot.MediaResponse.SignedUrl)
	return schemas.NewMediaFile(
		snapshot.Note,
		snapshot.Created,
		fmt.Sprintf("%s_%s", snapshot.FormatTimeStamp(), report.Dependent.DisplayName()),
		mediaUrl,
		snapshot.MediaResponse.MimeType,
	)
}

type Snapshot struct {
	Created       time.Time      `json:"created"`
	AttachmentId  string         `json:"attachment_id"`
	Note          string         `json:"note"`
	MediaResponse *MediaResponse `json:"media_response"`
}

type MediaResponse struct {
	SignedUrl string `json:"signed_url"`
	MimeType  string `json:"mime_type"`
}

func (s *Snapshot) FormatTimeStamp() string {
	return fmt.Sprintf("%d%02d%02d%02d%02d%02d",
		s.Created.Year(),
		s.Created.Month(),
		s.Created.Day(),
		s.Created.Hour(),
		s.Created.Minute(),
		s.Created.Second(),
	)
}

func (s *Snapshot) HydrateMediaData(
	client *http.Client,
	mediaUrl func(attachmentId string) *url.URL,
) error {
	data, err := fetchMediaData(client, mediaUrl(s.AttachmentId))
	if err != nil {
		return err
	}

	s.MediaResponse = data

	return nil
}

func fetchMediaData(client *http.Client, mediaUrl *url.URL) (media *MediaResponse, err error) {
	resp, err := client.Get(mediaUrl.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, utils.NewRequestError(resp, "Failed to fetch media data")
	}

	defer utils.CloseWithLog(resp.Body)
	body, _ := io.ReadAll(resp.Body)

	err = json.Unmarshal(body, &media)

	return media, err
}
