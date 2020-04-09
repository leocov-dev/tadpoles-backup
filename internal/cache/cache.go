package cache

import (
	"encoding/json"
	"github.com/leocov-dev/tadpoles-backup/config"
	"github.com/leocov-dev/tadpoles-backup/internal/api"
	"github.com/leocov-dev/tadpoles-backup/internal/utils"
	bolt "go.etcd.io/bbolt"
	"sort"
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

type ByEventTime []*api.Event

func (a ByEventTime) Len() int           { return len(a) }
func (a ByEventTime) Less(i, j int) bool { return a[i].EventTime.Time().Before(a[j].EventTime.Time()) }
func (a ByEventTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// return a list sorted by event time,
func ReadCache() (events []*api.Event, err error) {
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

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var event api.Event
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

	sort.Sort(ByEventTime(events))

	return events, nil
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
