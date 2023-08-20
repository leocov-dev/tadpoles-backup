package cache

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
	"path/filepath"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/api/tadpoles"
	"tadpoles-backup/internal/utils"
)

type TadpolesCache struct {
	bucketName string
	cookieFile string
	dbFile     string
}

func NewTadpolesCache() *TadpolesCache {
	cache := &TadpolesCache{
		bucketName: "TADPOLES_CACHE",
		cookieFile: filepath.Join(
			config.GetDataDir(),
			fmt.Sprintf(".%s-cookie", config.TADPOLES),
		),
		dbFile: filepath.Join(
			config.GetDataDir(),
			fmt.Sprintf(".%s-cache", config.TADPOLES),
		),
	}

	return cache
}

func (c *TadpolesCache) GetCookieFile() string {
	return c.cookieFile
}

func (c *TadpolesCache) ClearCache() error {
	return utils.DeleteFile(c.dbFile)
}

func (c *TadpolesCache) ClearLoginCookie() error {
	return utils.DeleteFile(c.cookieFile)
}

// ReadEventCache
// read the local bolt-db cache file and
// return a list of api events sorted by event time
func (c *TadpolesCache) ReadEventCache() (events tadpoles.Events, err error) {
	initializeBucket(c.dbFile, c.bucketName)

	db, err := bolt.Open(c.dbFile, 0600, &bolt.Options{ReadOnly: true})
	if err != nil {
		log.Errorln("failed bolt open")
		return nil, err
	}
	defer utils.CloseWithLog(db)

	log.Debug("reading event cache...")

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(c.bucketName))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var event tadpoles.Event
			err := json.Unmarshal(v, &event)
			if err != nil {
				return err
			}
			events = append(events, &event)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	events.Sort(func(e1, e2 *tadpoles.Event) bool {
		return e1.EventTime.Time().Before(e2.EventTime.Time())
	})

	return events, nil
}

// UpdateEventCache
// write a list of api events to the local bolt-db cache file
func (c *TadpolesCache) UpdateEventCache(events tadpoles.Events) error {
	initializeBucket(c.dbFile, c.bucketName)

	db, err := bolt.Open(c.dbFile, 0600, nil)
	if err != nil {
		return err
	}
	defer utils.CloseWithLog(db)

	log.Debug("writing event cache...")

	for _, event := range events {
		err := db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(c.bucketName))
			// TODO: maybe there is a way to store without marshaling
			//  to reduce processing on read-back
			j, err := json.Marshal(event)
			if err != nil {
				return err
			}
			return b.Put([]byte(event.EventKey), j)
		})
		if err != nil {
			return err
		}
	}

	return nil
}
