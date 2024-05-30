package tadpoles

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
	"path/filepath"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/cache"
	"tadpoles-backup/internal/utils"
)

type ApiCache struct {
	BucketName string
	CookieFile string
	DbFile     string
}

func NewApiCache(name string) *ApiCache {
	return &ApiCache{
		BucketName: fmt.Sprintf("TADPOLES_API_CACHE_%s", name),
		CookieFile: filepath.Join(
			config.GetDataDir(),
			fmt.Sprintf(".%s-api-cookie", name),
		),
		DbFile: filepath.Join(
			config.GetDataDir(),
			fmt.Sprintf(".%s-api-c", name),
		),
	}
}

func (c *ApiCache) ClearCache() error {
	return utils.DeleteFile(c.DbFile)
}

func (c *ApiCache) ClearLoginCookie() error {
	return utils.DeleteFile(c.CookieFile)
}

// ReadEventCache
// read the local bolt-db c file and
// return a list of api Events sorted by Event time
func (c *ApiCache) readEventCache() (events Events, err error) {
	cache.InitializeBucket(c.DbFile, c.BucketName)

	db, err := bolt.Open(c.DbFile, 0600, &bolt.Options{ReadOnly: true})
	if err != nil {
		log.Errorln("failed bolt open")
		return nil, err
	}
	defer utils.CloseWithLog(db)

	log.Debug("reading Event cache...")

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(c.BucketName))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var event Event
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

	events.Sort(func(e1, e2 *Event) bool {
		return e1.EventTime.Time().Before(e2.EventTime.Time())
	})

	log.Debug(fmt.Sprintf("cached Events: %d", len(events)))

	return events, nil
}

// UpdateEventCache
// write a list of api Events to the local bolt-db c file
func (c *ApiCache) updateEventCache(events Events) error {
	cache.InitializeBucket(c.DbFile, c.BucketName)

	db, err := bolt.Open(c.DbFile, 0600, nil)
	if err != nil {
		return err
	}
	defer utils.CloseWithLog(db)

	log.Debug("writing Event cache...")

	for _, event := range events {
		err := db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(c.BucketName))
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
