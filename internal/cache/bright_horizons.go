package cache

import (
	"bufio"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
	"os"
	"path/filepath"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/api/bright_horizons"
	"tadpoles-backup/internal/utils"
)

type BrightHorizonsCache struct {
	bucketPrefix string
	dbFile       string
	apiKeyFile   string
}

func NewBrightHorizonsCache() *BrightHorizonsCache {
	return &BrightHorizonsCache{
		bucketPrefix: "BH_",
		dbFile: filepath.Join(
			config.GetDataDir(),
			fmt.Sprintf(".%s-cache", config.BRIGHT_HORIZONS),
		),
		apiKeyFile: filepath.Join(
			config.GetDataDir(),
			fmt.Sprintf(".%s-key", config.BRIGHT_HORIZONS),
		),
	}
}

func (c *BrightHorizonsCache) GetApiKey() string {
	if utils.FileExists(c.apiKeyFile) {
		f, err := os.Open(c.apiKeyFile)
		if err != nil {
			return ""
		}
		defer utils.CloseWithLog(f)

		fileScanner := bufio.NewScanner(f)

		fileScanner.Scan()
		return fileScanner.Text()
	}

	return ""
}

func (c *BrightHorizonsCache) StoreApiKey(apiKey string) error {
	err := os.WriteFile(c.apiKeyFile, []byte(apiKey), 0600)
	if err != nil {
		log.Debug("Failed to write cookies json to file...", err)
		return err
	}

	return nil
}

func (c *BrightHorizonsCache) dependentBucketName(dependentId string) string {
	return fmt.Sprintf("%s_%s", c.bucketPrefix, dependentId)
}

func (c *BrightHorizonsCache) ReadReportCache(dependentId string) (reports bright_horizons.Reports, err error) {
	initializeBucket(c.dbFile, c.dependentBucketName(dependentId))

	db, err := bolt.Open(c.dbFile, 0600, &bolt.Options{ReadOnly: true})
	if err != nil {
		log.Errorln("failed bolt open")
		return nil, err
	}
	defer utils.CloseWithLog(db)

	log.Debug("reading report cache...")

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(c.dependentBucketName(dependentId)))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var report bright_horizons.Report
			err := json.Unmarshal(v, &report)
			if err != nil {
				return err
			}
			reports = append(reports, &report)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	reports.Sort(bright_horizons.ByReportDate)

	return reports, nil
}

func (c *BrightHorizonsCache) UpdateReportCache(
	dependentId string,
	reports bright_horizons.Reports,
) error {
	initializeBucket(c.dbFile, c.dependentBucketName(dependentId))

	db, err := bolt.Open(c.dbFile, 0600, nil)
	if err != nil {
		return err
	}
	defer utils.CloseWithLog(db)

	log.Debugf("writing event cache %s <- %d", dependentId, len(reports))

	for _, report := range reports {
		updateErr := db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(c.dependentBucketName(dependentId)))
			// TODO: maybe there is a way to store without marshaling
			//  to reduce processing on read-back
			j, jsonErr := json.Marshal(report)
			if jsonErr != nil {
				return jsonErr
			}
			return b.Put([]byte(report.Id), j)
		})
		if updateErr != nil {
			return updateErr
		}
	}

	return nil
}

func (c *BrightHorizonsCache) ClearCache() error {
	return utils.DeleteFile(c.dbFile)
}

func (c *BrightHorizonsCache) ClearLoginData() error {
	return utils.DeleteFile(c.apiKeyFile)
}
