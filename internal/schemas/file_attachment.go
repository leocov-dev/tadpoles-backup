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
	log "github.com/sirupsen/logrus"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/api"
	"tadpoles-backup/internal/utils"
	"time"
)

type FileAttachment struct {
	Comment       string
	AttachmentKey string
	EventKey      string
	ChildName     string
	EventTime     time.Time
	tempFile      string
	imageType     types.Type
	EventMime     string
}

func NewFileAttachment(event *api.Event, eventAttachment *api.Attachment) *FileAttachment {
	return &FileAttachment{
		Comment:       event.Comment,
		AttachmentKey: eventAttachment.AttachmentKey,
		EventKey:      event.EventKey,
		ChildName:     event.ChildName,
		EventTime:     event.EventTime.Time(),
		EventMime:     eventAttachment.MimeType,
	}
}

// get the save target file name without extension
func (a *FileAttachment) SaveName() string {
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
func (a *FileAttachment) saveTarget(backupRoot string) (filePath string, err error) {
	if a.imageType.Extension == "" {
		return "", errors.New("must call Download() in order to establish file extension")
	}
	dir := filepath.Join(backupRoot, fmt.Sprint(a.EventTime.Year()), fmt.Sprintf("%d-%02d-%02d", a.EventTime.Year(), a.EventTime.Month(), a.EventTime.Day()))
	fileName := fmt.Sprintf("%s.%s", a.SaveName(), a.imageType.Extension)
	return filepath.Join(dir, fileName), nil
}

// Download
// fetch file from url to a temporary directory
func (a *FileAttachment) Download() (err error) {
	log.Debug("Downloading: ", a.AttachmentKey)
	log.Debugf("%s %s %s\n", a.ChildName, a.Comment, a.EventTime)

	resp, err := api.S.GetAttachment(a.EventKey, a.AttachmentKey)
	if err != nil {
		return err
	}
	defer utils.CloseWithLog(resp.Body)

	tempFile, err := os.CreateTemp(config.TempDir, config.TempFilePattern)
	if err != nil {
		return err
	}
	defer utils.CloseWithLog(tempFile)

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return nil
	}

	a.tempFile = tempFile.Name()

	a.imageType, err = filetype.MatchFile(a.tempFile)
	if err != nil {
		return err
	}

	return nil
}

// Save
// create the necessary directories and move the temporary file to the target with a new name
func (a *FileAttachment) Save(backupRoot string) (err error) {
	if a.tempFile == "" {
		return errors.New("must call Download() before Save()")
	}

	if utils.IsImageType(a.imageType) {
		err = a.processImageFile()
		if err != nil {
			return err
		}
	}

	savePath, err := a.saveTarget(backupRoot)
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

	if a.imageType != matchers.TypeJpeg {
		// handle Non-jpeg files
		log.Debugf("Not jpg, converting: %s\n\n", a.imageType)

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

	switch a.imageType {
	case matchers.TypePng:
		img, err = png.Decode(tempFile)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New(fmt.Sprintf("jpeg conversion not implemented for %s", a.imageType.Extension))
	}

	jpgBuffer := new(bytes.Buffer)

	err = jpeg.Encode(jpgBuffer, img, &jpeg.Options{Quality: 85})
	if err != nil {
		return nil, err
	}

	a.imageType = matchers.TypeJpeg

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

type FileAttachments []*FileAttachment

type FileAttachmentMap map[string]FileAttachments

func (fam FileAttachmentMap) PrettyPrint(heading string) {
	utils.WriteMain(heading, "")
	for k, v := range fam {
		utils.WriteSub(k, fmt.Sprint(len(v)))
	}
}
