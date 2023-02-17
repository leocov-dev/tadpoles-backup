package cache

import (
	"encoding/json"
	bolt "go.etcd.io/bbolt"
	"sort"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/api"
	"tadpoles-backup/internal/utils"
)

var (
	eventsBucket = []byte("EVENTS")
)

func OpenCacheDB(options *bolt.Options) (*bolt.DB, error) {
	return bolt.Open(config.GetCacheDbFile(), 0600, options)
}

func InitializeCache() error {
	db, err := OpenCacheDB(nil)
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

// ReadEventCache
// read the local bolt-db cache file and
// return a list of api events sorted by event time
func ReadEventCache() (events []*api.Event, err error) {
	db, err := OpenCacheDB(&bolt.Options{ReadOnly: true})
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

// WriteEventCache
// write a list of api events to the local bolt-db cache file
func WriteEventCache(events []*api.Event) error {
	db, err := OpenCacheDB(nil)
	if err != nil {
		return err
	}
	defer utils.CloseWithLog(db)

	for _, event := range events {
		err := db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(eventsBucket)
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
