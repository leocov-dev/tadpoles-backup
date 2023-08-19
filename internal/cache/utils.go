package cache

import bolt "go.etcd.io/bbolt"

func initializeBucket(dbFile, bucketName string) {
	db, _ := bolt.Open(dbFile, 0600, nil)
	_ = db.Update(func(tx *bolt.Tx) error {
		_, _ = tx.CreateBucketIfNotExists([]byte(bucketName))
		return nil
	})
	_ = db.Close()
}
