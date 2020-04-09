package cache

import (
	"encoding/json"
	"github.com/leocov-dev/tadpoles-backup/config"
	"github.com/leocov-dev/tadpoles-backup/internal/api"
	"github.com/leocov-dev/tadpoles-backup/internal/utils"
	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

var (
	eventsBucket = []byte("EVENTS")
)

func InitializeCache() error {
	db, err := bolt.Open(config.TadpolesCacheFile, 0600, nil)
	if err != nil {
		return err
	}
	defer utils.CloseWithLog(db)

	return db.Update(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		_, err := tx.CreateBucketIfNotExists(eventsBucket)
		if err != nil {
			return err
		}

		return nil
	})
}

func Read() (*api.Event, error) {
	db, err := bolt.Open(config.TadpolesCacheFile, 0600,
		&bolt.Options{
			ReadOnly: true,
		})
	if err != nil {
		return nil, err
	}
	defer utils.CloseWithLog(db)

	err = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(eventsBucket)

		c := b.Cursor()

		k, v := c.First()
		log.Debugf("key=%s, value=%s\n", k, v)
		k, v = c.Next()
		log.Debugf("key=%s, value=%s\n", k, v)
		k, v = c.Next()
		log.Debugf("key=%s, value=%s\n", k, v)
		k, v = c.Next()
		log.Debugf("key=%s, value=%s\n", k, v)
		k, v = c.Next()
		log.Debugf("key=%s, value=%s\n", k, v)
		k, v = c.Next()
		log.Debugf("key=%s, value=%s\n", k, v)
		k, v = c.Next()
		log.Debugf("key=%s, value=%s\n", k, v)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func StoreEvents(events []*api.Event) error {
	db, err := bolt.Open(config.TadpolesCacheFile, 0600, nil)
	if err != nil {
		return err
	}
	defer utils.CloseWithLog(db)

	for _, event := range events {
		err := db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(eventsBucket)
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
