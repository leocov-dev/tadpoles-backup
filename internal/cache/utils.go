package cache

import (
	bolt "go.etcd.io/bbolt"
	"tadpoles-backup/internal/utils"
)

func InitializeBucket(dbFile, bucketName string) {
	db, _ := bolt.Open(dbFile, 0600, nil)
	defer utils.CloseWithLog(db)
	_ = db.Update(func(tx *bolt.Tx) error {
		_, _ = tx.CreateBucketIfNotExists([]byte(bucketName))
		return nil
	})
}

func DeleteBucket(dbFile, bucketName string) {
	db, _ := bolt.Open(dbFile, 0600, nil)
	defer utils.CloseWithLog(db)
	_ = db.Update(func(tx *bolt.Tx) error {
		if b := tx.Bucket([]byte(bucketName)); b != nil {
			_ = b.DeleteBucket([]byte(bucketName))
		}
		return nil
	})
}
