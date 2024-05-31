package http_utils

import (
	"context"
	"fmt"
	"github.com/h2non/filetype/types"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/interfaces"
	"tadpoles-backup/internal/utils"
	"tadpoles-backup/internal/utils/progress"
	"tadpoles-backup/pkg/async"
	"time"
)

type MediaFile struct {
	comment     string
	timestamp   time.Time
	downloadUrl *url.URL
	FileName    string
	MimeType    types.MIME
}

func NewMediaFile(
	comment string,
	timestamp time.Time,
	fileName string,
	downloadUrl *url.URL,
	mimeType string,
) MediaFile {
	return MediaFile{
		comment:     comment,
		timestamp:   timestamp,
		FileName:    fileName,
		downloadUrl: downloadUrl,
		MimeType:    types.NewMIME(mimeType),
	}
}

func (m MediaFile) FilePath(rootDir string) string {
	// final file extension will be determined at download time when parsing
	// actual data
	return filepath.Join(
		rootDir,
		fmt.Sprint(m.timestamp.Year()),
		fmt.Sprintf("%d-%02d-%02d",
			m.timestamp.Year(), m.timestamp.Month(), m.timestamp.Day()),
		fmt.Sprintf("%s.partial", m.FileName),
	)
}

func (m MediaFile) GetTimestamp() time.Time {
	return m.timestamp
}

func (m MediaFile) GetComment() string {
	return m.comment
}

func (m MediaFile) Download(client interfaces.HttpClient, dlRoot string) error {
	return DownloadFile(
		client,
		m.downloadUrl,
		m.FilePath(dlRoot),
		m,
	)
}

type MediaFiles []MediaFile

type MediaFileCountMap map[string]int

func (mfs MediaFiles) CountByType() MediaFileCountMap {
	countMap := make(MediaFileCountMap)
	countMap["Images"] = 0
	countMap["Videos"] = 0
	countMap["Unknown"] = 0

	for _, mediaFile := range mfs {
		if utils.IsImageType(mediaFile.MimeType) {
			countMap["Images"] += 1
		} else if utils.IsVideoType(mediaFile.MimeType) {
			countMap["Videos"] += 1
		} else {
			countMap["Unknown"] += 1
		}
	}

	return countMap
}

func (mfs MediaFiles) FilterOnlyNew(
	downloadRoot string,
) (onlyNew MediaFiles, err error) {
	filePathMap := make(map[string]MediaFile)

	for _, f := range mfs {
		filePathMap[f.FileName] = f
	}

	err = filepath.Walk(downloadRoot,
		func(path_ string, info_ os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info_.IsDir() {
				return nil
			}

			// if the path_ matches one in the map
			// then it is already downloaded, remove it
			check := strings.TrimSuffix(filepath.Base(path_), filepath.Ext(path_))
			if _, ok := filePathMap[check]; ok {
				delete(filePathMap, check)
			}
			return nil
		})

	// construct a slice of new (not downloaded) files
	for _, v := range filePathMap {
		onlyNew = append(onlyNew, v)
	}

	return onlyNew, err
}

func (mfs MediaFiles) DownloadAll(
	client interfaces.HttpClient,
	dlRoot string,
	ctx context.Context,
	sharedProgressBar *progress.BarWrapper,
) error {
	taskPool := async.NewTaskPool(ctx, nil)

	for _, f := range mfs {
		select {
		case <-ctx.Done():
			return async.NewCanceledError()
		default:
			err := taskPool.AddTask(
				NewDownloadTask(
					client,
					f,
					dlRoot,
					sharedProgressBar,
				),
			)
			if err != nil {
				return err
			}
		}
	}
	taskPool.Wait()

	return taskPool.Errors()
}

func (cm MediaFileCountMap) PrettyPrint(heading string) {
	utils.WriteMain(heading, "")
	for k, v := range cm {
		if !config.DebugMode && k == "Unknown" {
			continue
		}
		utils.WriteSub(k, fmt.Sprint(v))
	}
}
