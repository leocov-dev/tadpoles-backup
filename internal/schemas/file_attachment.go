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
	if utils.IsImageType(a.ImageType) {
		err = a.convertToJpgIfRequired()
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

// TODO: this is messy
func (a *FileAttachment) convertToJpgIfRequired() (err error) {
	if a.ImageType == matchers.TypeJpeg {
		log.Debug("Already jpg...\n")

		err = writeExifTag(a)
		if err != nil {
			return err
		}

		return nil
	}

	log.Debugf("Not jpg, converting: %s\n\n", a.ImageType)

	tempFile, err := os.OpenFile(a.tempFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer utils.CloseWithLog(tempFile)

	var img image.Image

	switch a.ImageType {
	case matchers.TypePng:
		img, err = png.Decode(tempFile)
		if err != nil {
			return err
		}

	default:
		return errors.New(fmt.Sprintf("jpeg conversion not implemented for %s", a.ImageType.Extension))
	}

	jpgBytes := new(bytes.Buffer)

	err = jpeg.Encode(jpgBytes, img, &jpeg.Options{Quality: 85})
	if err != nil {
		return err
	}
	a.ImageType = matchers.TypeJpeg

	jmp := gjis.NewJpegMediaParser()
	sl, err := jmp.ParseBytes(jpgBytes.Bytes())
	if err != nil {
		return err
	}

	err = setExifData(sl, a.EventTime, a.Comment)
	if err != nil {
		return err
	}

	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	err = tempFile.Truncate(0)
	if err != nil {
		return err
	}
	err = sl.Write(tempFile)
	if err != nil {
		log.Debug("Failed sl.Write()...\n")
		return err
	}

	return nil
}

func writeExifTag(attachment *FileAttachment) (err error) {
	log.Debug("writeExifTag()...\n")
	jmp := gjis.NewJpegMediaParser()

	sl, err := jmp.ParseFile(attachment.tempFile)
	if err != nil {
		return err
	}

	err = setExifData(sl, attachment.EventTime, attachment.Comment)
	if err != nil {
		return err
	}

	f, err := os.Create(attachment.tempFile)
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
