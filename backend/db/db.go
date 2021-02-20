package db

import (
	"fmt"

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

func (db *DB) GetHash(dir string, filename string) ([]byte, error) {
	var hash []byte
	err := db.View(func(tx *bolt.Tx) error {
		bucketFolders := tx.Bucket([]byte(key))
		bucketHashes := bucketFolders.Bucket([]byte(dir))
		hash = bucketHashes.Get([]byte(filename))
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("view: %w", err)
	}
	return hash, nil
}

func (db *DB) SetHash(dir string, filename string, hash []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucketFolders := tx.Bucket([]byte(key))
		bucketHashes, err := bucketFolders.CreateBucketIfNotExists([]byte(dir))
		if err != nil {
			return fmt.Errorf("create dir bucket: %w", err)
		}
		if err := bucketHashes.Put([]byte(filename), hash); err != nil {
			return fmt.Errorf("insert mod time: %w", err)
		}
		return nil
	})
}
