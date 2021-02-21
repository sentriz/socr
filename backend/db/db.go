package db

import (
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

const (
	key = "folders"
)

type DB struct {
	*bolt.DB
}

func New(path string) (*DB, error) {
	bdb, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}
	err = bdb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(key))
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("create initial bucket: %w", err)
	}
	return &DB{bdb}, nil
}

func (db *DB) GetModTime(dir string, filename string) (*time.Time, error) {
	var modTime *time.Time
	err := db.View(func(tx *bolt.Tx) error {
		bucketFolders := tx.Bucket([]byte(key))
		bucketModTimes := bucketFolders.Bucket([]byte(dir))
		modTimeRaw := bucketModTimes.Get([]byte(filename))
		if modTimeRaw == nil {
			return nil
		}
		modTime = &time.Time{}
		if err := modTime.UnmarshalBinary(modTimeRaw); err != nil {
			return fmt.Errorf("decode time: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("view: %w", err)
	}
	return modTime, nil
}

func (db *DB) SetModTime(dir string, filename string, modTime time.Time) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucketFolders := tx.Bucket([]byte(key))
		bucketModTimees, err := bucketFolders.CreateBucketIfNotExists([]byte(dir))
		if err != nil {
			return fmt.Errorf("create dir bucket: %w", err)
		}
		modTimeRaw, err := modTime.MarshalBinary()
		if err != nil {
			return fmt.Errorf("encode time: %w", err)
		}
		if err := bucketModTimees.Put([]byte(filename), modTimeRaw); err != nil {
			return fmt.Errorf("insert mod time: %w", err)
		}
		return nil
	})
}
