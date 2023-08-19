package schemas

import (
	"context"
	"fmt"
	"github.com/h2non/filetype/types"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/utils"
	"tadpoles-backup/internal/utils/progress"
	"time"
)

type MediaFile struct {
	comment     string
	timestamp   time.Time
	FileName    string
	DownloadUrl *url.URL
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
		DownloadUrl: downloadUrl,
		MimeType:    types.NewMIME(mimeType),
	}
}

func (f MediaFile) FilePath(rootDir string) string {
	// final file extension will be determined at download time when parsing
	// actual data
	return filepath.Join(
		rootDir,
		fmt.Sprint(f.timestamp.Year()),
		fmt.Sprintf("%d-%02d-%02d",
			f.timestamp.Year(), f.timestamp.Month(), f.timestamp.Day()),
		fmt.Sprintf("%s.partial", f.FileName),
	)
}

func (f MediaFile) GetTimestamp() time.Time {
	return f.timestamp
}

func (f MediaFile) GetComment() string {
	return f.comment
}

func (f MediaFile) Download(client *http.Client, dlRoot string) error {
	return utils.DownloadFile(
		client,
		f.DownloadUrl,
		f.FilePath(dlRoot),
		f,
	)
}

type MediaFiles []MediaFile

type MediaFileCountMap map[string]int

func (fs MediaFiles) CountByType() MediaFileCountMap {
	countMap := make(MediaFileCountMap)
	countMap["Images"] = 0
	countMap["Videos"] = 0
	countMap["Unknown"] = 0

	for _, f := range fs {
		if utils.IsImageType(f.MimeType) {
			countMap["Images"] += 1
		} else if utils.IsVideoType(f.MimeType) {
			countMap["Videos"] += 1
		} else {
			countMap["Unknown"] += 1
		}
	}

	return countMap
}

func (fs MediaFiles) FilterOnlyNew(
	downloadRoot string,
) (onlyNew MediaFiles, err error) {
	filePathMap := make(map[string]MediaFile)

	for _, f := range fs {
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

func (fs MediaFiles) DownloadAll(
	client *http.Client,
	dlRoot string,
	concurrency int,
	ctx context.Context,
	sharedProgressBar *progress.BarWrapper,
) (dlErrors []string) {
	errorChan := make(chan string)

	downloadPool := NewDownloadPool(concurrency)

	for _, f := range fs {
		downloadPool.Add(
			NewDownloadTask(
				client,
				f,
				dlRoot,
				errorChan,
				ctx,
				sharedProgressBar,
			),
		)
	}
	downloadPool.ProcessTasks()
	close(errorChan)

	for s := range errorChan {
		dlErrors = append(dlErrors, s)
	}

	return dlErrors
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
