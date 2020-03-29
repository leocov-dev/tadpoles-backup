package schemas

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
	"github.com/leocov-dev/tadpoles-backup/config"
	"github.com/leocov-dev/tadpoles-backup/internal/api"
	"github.com/leocov-dev/tadpoles-backup/internal/utils"
	log "github.com/sirupsen/logrus"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type FileAttachment struct {
	Comment           string
	AttachmentKey     string
	EventKey          string
	ChildName         string
	CreateTime        time.Time
	EventTime         time.Time
	tempFile          string
	AlreadyDownloaded bool
	ImageType         types.Type
	EventMime         string
}

// get the save target file name without extension
func (a *FileAttachment) GetSaveName() string {
	timestamp := fmt.Sprintf("%d%02d%02d%02d%02d%02d",
		a.EventTime.Year(),
		a.EventTime.Month(),
		a.EventTime.Day(),
		a.EventTime.Hour(),
		a.EventTime.Minute(),
		a.EventTime.Second(),
	)
	return fmt.Sprintf("%s_%s", timestamp, a.ChildName)
}

// get the path and filename for the final save location
func (a *FileAttachment) GetSaveTarget(backupRoot string) (filePath string, err error) {
	dir := filepath.Join(backupRoot, fmt.Sprint(a.EventTime.Year()), fmt.Sprintf("%d-%02d-%02d", a.EventTime.Year(), a.EventTime.Month(), a.EventTime.Day()))
	fileName := fmt.Sprintf("%s.%s", a.GetSaveName(), a.ImageType.Extension)
	return filepath.Join(dir, fileName), nil
}

// download file to a temporary directory
func (a *FileAttachment) Download() (err error) {
	log.Debug("Downloading: ", a.AttachmentKey)
	log.Debugf("%s %s %s\n", a.ChildName, a.Comment, a.EventTime)

	resp, err := api.GetAttachment(a.EventKey, a.AttachmentKey)
	if err != nil {
		return err
	}
	defer utils.CloseWithLog(resp.Body)

	tempFile, err := ioutil.TempFile(config.TempDir, config.TempFilePattern)
	if err != nil {
		return err
	}
	defer utils.CloseWithLog(tempFile)

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return nil
	}

	a.tempFile = tempFile.Name()

	a.ImageType, err = filetype.MatchFile(a.tempFile)
	if err != nil {
		return err
	}

	return nil
}

// create the necessary directories and move the temporary file to the target with a new name
func (a *FileAttachment) Save(backupRoot string) (err error) {
	if a.tempFile == "" {
		return errors.New("must call Download() before Save()")
	}

	if utils.IsImageType(a.ImageType) {
		err = a.processImageFile()
		if err != nil {
			return err
		}
	}

	savePath, err := a.GetSaveTarget(backupRoot)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(savePath), os.ModePerm)
	if err != nil {
		return err
	}

	log.Debugf("Saving to: %s\n\n", savePath)

	time.Sleep(2 * time.Second)
	err = utils.CopyFile(a.tempFile, savePath)
	if err != nil {
		return err
	}

	return nil
}

func (a *FileAttachment) processImageFile() (err error) {
	jmp := gjis.NewJpegMediaParser()
	var sl *gjis.SegmentList

	if a.ImageType != matchers.TypeJpeg {
		// handle Non-jpeg files
		log.Debugf("Not jpg, converting: %s\n\n", a.ImageType)

		jpegBytes, err := a.convertToJpeg()
		if err != nil {
			return err
		}

		sl, err = jmp.ParseBytes(jpegBytes)
		if err != nil {
			return err
		}
	} else {
		// handle jpeg files
		sl, err = jmp.ParseFile(a.tempFile)
		if err != nil {
			return err
		}
	}

	err = a.writeExifTag(sl)
	if err != nil {
		return err
	}

	return nil
}

// decode images from other formats to raw image object
func (a *FileAttachment) convertToJpeg() (jpegBytes []byte, err error) {
	tempFile, err := os.OpenFile(a.tempFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer utils.CloseWithLog(tempFile)
	if err != nil {
		return nil, err
	}

	var img image.Image

	switch a.ImageType {
	case matchers.TypePng:
		img, err = png.Decode(tempFile)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New(fmt.Sprintf("jpeg conversion not implemented for %s", a.ImageType.Extension))
	}

	jpgBuffer := new(bytes.Buffer)

	err = jpeg.Encode(jpgBuffer, img, &jpeg.Options{Quality: 85})
	if err != nil {
		return nil, err
	}

	a.ImageType = matchers.TypeJpeg

	return jpgBuffer.Bytes(), nil
}

func (a *FileAttachment) writeExifTag(sl *gjis.SegmentList) (err error) {
	err = setExifData(sl, a.EventTime, a.Comment)
	if err != nil {
		return err
	}

	f, err := os.Create(a.tempFile)
	if err != nil {
		return err
	}
	defer utils.CloseWithLog(f)

	err = sl.Write(f)
	if err != nil {
		return err
	}

	return nil
}

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
