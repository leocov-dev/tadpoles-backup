package cache

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"tadpoles-backup/config"
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
		logrus.Debug("Failed to write cookies json to file...", err)
		return err
	}

	return nil
}

func (c *BrightHorizonsCache) dependentBucketName(dependentId string) []byte {
	return []byte(fmt.Sprintf("%s_%s", c.bucketPrefix, dependentId))
}

func (c *BrightHorizonsCache) ClearCache() error {
	return utils.DeleteFile(c.dbFile)
}

func (c *BrightHorizonsCache) ClearLoginData() error {
	return utils.DeleteFile(c.apiKeyFile)
}
