package http_utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dsoprea/go-exif/v2"
	exifcommon "github.com/dsoprea/go-exif/v2/common"
	gjis "github.com/dsoprea/go-jpeg-image-structure"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/matchers"
	"github.com/h2non/filetype/types"
	log "github.com/sirupsen/logrus"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"tadpoles-backup/internal/interfaces"
	"tadpoles-backup/internal/utils"
	"time"
)

type MediaMetadata interface {
	GetTimestamp() time.Time
	GetComment() string
}

func DownloadFile(
	client interfaces.HttpClient,
	fileUrl *url.URL,
	targetPath string,
	metadata MediaMetadata,
) error {
	log.Debug("Download To:", targetPath)

	resp, err := client.Get(fileUrl.String())
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return utils.NewRequestError(resp, "could not get attachment")
	}
	defer utils.CloseWithLog(resp.Body)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fileType, err := filetype.Match(data)
	log.Debug("type: ", fileType)

	buffer := bytes.NewBuffer(data)
	if utils.IsImageType(fileType.MIME) {
		// If the file data is an image (jpeg, png, etc.)
		// then we will convert it to a jpeg if required
		// and apply Exif tag data.
		if fileType != matchers.TypeJpeg {
			log.Debug("Not jpg, converting: ", fileType)
			err := convertToJpeg(buffer, fileType)
			if err != nil {
				return err
			}
			fileType = matchers.TypeJpeg
		}

		err = writeExifTag(buffer, metadata.GetTimestamp(), metadata.GetComment())
		if err != nil {
			return err
		}
		data = buffer.Bytes()

	} else if utils.IsVideoType(fileType.MIME) {
		// TODO: video tags
	} else {
		return nil
	}

	// create the parent directories for the file
	err = os.MkdirAll(filepath.Dir(targetPath), os.ModePerm)
	if err != nil {
		return err
	}

	// Create a file with the correct extension now that we've parsed it
	// directly from the data header
	targetFile, err := os.Create(fmt.Sprintf("%s.%s",
		strings.TrimSuffix(targetPath, filepath.Ext(targetPath)),
		fileType.Extension,
	))
	if err != nil {
		return err
	}
	defer utils.CloseWithLog(targetFile)

	_, err = targetFile.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func convertToJpeg(buffer *bytes.Buffer, fileType types.Type) error {
	var (
		img image.Image
		err error
	)

	switch fileType {
	case matchers.TypePng:
		log.Debug("start png decode")
		img, err = png.Decode(buffer)
		if err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("jpeg conversion not implemented for %s", fileType.Extension))
	}

	log.Debug("start jpeg encode")
	buffer.Reset()
	err = jpeg.Encode(buffer, img, &jpeg.Options{Quality: 85})
	if err != nil {
		return err
	}

	return nil
}

func writeExifTag(buffer *bytes.Buffer, timestamp time.Time, comment string) (err error) {
	jmp := gjis.NewJpegMediaParser()
	var sl *gjis.SegmentList

	sl, err = jmp.Parse(buffer, buffer.Len())
	if err != nil {
		return err
	}

	err = setExifData(sl, timestamp, comment)
	if err != nil {
		return err
	}

	buffer.Reset()
	err = sl.Write(buffer)
	if err != nil {
		return err
	}

	return nil
}

// setExifData
// https://github.com/dsoprea/go-jpeg-image-structure/blob/master/jpeg_test.go
func setExifData(sl *gjis.SegmentList, dateTime time.Time, userComment string) (err error) {
	log.Debug("setExifData()...\n")

	rootIb, err := sl.ConstructExifBuilder()
	if err != nil {
		im := exif.NewIfdMappingWithStandard()
		ti := exif.NewTagIndex()
		err := exif.LoadStandardTags(ti)
		if err != nil {
			return err
		}

		rootIb = exif.NewIfdBuilder(im, ti, exifcommon.IfdPathStandard, exifcommon.EncodeDefaultByteOrder)
	}

	ifdIb, err := exif.GetOrCreateIbFromRootIb(rootIb, "IFD0")
	if err != nil {
		log.Debug("Failed exif.GetOrCreateIbFromRootIb()...\n")
		return err
	}

	// DateTime
	updatedTimestampPhrase := exif.ExifFullTimestampString(dateTime)
	err = ifdIb.SetStandardWithName("DateTime", updatedTimestampPhrase)
	if err != nil {
		return err
	}

	// ImageDescription
	err = ifdIb.SetStandardWithName("ImageDescription", userComment)
	if err != nil {
		return err
	}

	// Update the exif segment.
	err = sl.SetExif(rootIb)
	if err != nil {
		log.Debug("Failed sl.SetExif()...\n")
		return err
	}

	return nil
}
