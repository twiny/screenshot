package db

import (
	"encoding/json"
	"errors"
	"sync/atomic"
	"time"

	"github.com/twiny/screenshot/cmd/screen/api"
	"go.etcd.io/bbolt"
)

var (
	ErrImageNotFound = errors.New("image not found")
)

// DB
type Store struct {
	db *bbolt.DB
}

// NewStore
func NewStore(path string) (*Store, error) {
	db, err := bbolt.Open(path, 0644, bbolt.DefaultOptions)
	if err != nil {
		return nil, err
	}

	// create bucket
	if err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("images"))
		return err
	}); err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

// SaveImage
func (s *Store) SaveImage(img *api.Image) error {
	return s.db.Batch(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("images"))

		// encode
		data, err := json.Marshal(img)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(img.UUID), data)
	})
}

// FindImage
func (s *Store) FindImage(key string) (*api.Image, error) {
	img := api.Image{}
	if err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("images"))

		data := bucket.Get([]byte(key))
		if data == nil {
			return ErrImageNotFound
		}
		return json.Unmarshal(data, &img)
	}); err != nil {
		return nil, err
	}

	return &img, nil
}

// Clean
func (s *Store) Clean(d time.Duration) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("images"))

		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			img := api.Image{}
			if err := json.Unmarshal(v, &img); err != nil {
				return err
			}

			if time.Since(img.CreatedAt) > d {
				if err := bucket.Delete(k); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// Statistic
func (s *Store) Statistic() (map[string]interface{}, error) {
	stats := map[string]interface{}{}
	var success, fail, pending int32 = 0, 0, 0

	if err := s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("images"))

		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			img := api.Image{}
			if err := json.Unmarshal(v, &img); err != nil {
				return err
			}
			// count
			switch img.Status {
			case api.ImageStatusSuccess:
				atomic.AddInt32(&success, 1)
			case api.ImageStatusFail:
				atomic.AddInt32(&fail, 1)
			case api.ImageStatusPending:
				atomic.AddInt32(&pending, 1)
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	//
	stats["success"] = success
	stats["fail"] = fail
	stats["pending"] = pending

	return stats, nil

}

// Close
func (s *Store) Close() {
	s.db.Close()
}
