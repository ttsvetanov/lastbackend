package db

import (
	"errors"
	"github.com/boltdb/bolt"
	"github.com/deployithq/deployit/drivers/interfaces"
)

const base = "map"

type Bolt struct {
	DB *bolt.DB
}

func Open(log interfaces.Log, path string) *bolt.DB {

	db, err := bolt.Open(path, 0766, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(base))
		if err != nil {
			log.Error(err)
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return db

}

func (b *Bolt) Write(log interfaces.Log, key, value string) error {
	log.Debug("Add file to storage")

	err := b.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(base))
		if err != nil {
			log.Error(err)
			return err
		}

		err = bucket.Put([]byte(key), []byte(value))
		if err != nil {
			log.Error(err)
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (b *Bolt) Read(log interfaces.Log, key string) (string, error) {
	log.Debug("Read file from storage")

	var val string

	err := b.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(base))
		if bucket == nil {
			err := errors.New("BUCKET_NOT_FOUND")
			log.Error(err)
			return err
		}

		val = string(bucket.Get([]byte(key)))

		return nil
	})

	if err != nil {
		return val, err
	}

	return val, nil
}

func (b *Bolt) Delete(log interfaces.Log, key string) error {
	log.Debug("Delete file from storage")

	err := b.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(base))
		if bucket == nil {
			log.Error(interfaces.ErrBucketNotFound)
			return interfaces.ErrBucketNotFound
		}

		err := bucket.Delete([]byte(key))
		if err != nil {
			log.Error(err)
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (b *Bolt) ListAllFiles(log interfaces.Log) ([]string, error) {
	log.Debug("List all files from storage")

	var files []string

	err := b.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(base))

		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			if k != nil {
				files = append(files, string(k))
			} else {
				break
			}
		}

		return nil
	})

	if err != nil {
		log.Error(err)
		return files, err
	}

	return files, nil
}
