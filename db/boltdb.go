package db

import (
	"time"

	"github.com/boltdb/bolt"
)

type Interface interface {
	CreateBucket(name string) error
	DeleteBucket(name string) error
	Buckets() ([]string, error)
	UseBucket(name string) error
	CurrentBucket() string

	Put(key, value string) error
	Remove(key string) error
	Get(key string) (string, error)
	Keys() ([]string, error)
	List() (map[string]string, error)

	Close()
}

type boltDB struct {
	db *bolt.DB

	currBucket []byte
}

func New(path string) (Interface, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		return nil, err
	}

	return &boltDB{
		db: db,
	}, nil
}

func (db *boltDB) Close() {
	if db.db != nil {
		db.db.Close()
	}
}

func (db *boltDB) CreateBucket(name string) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(name))
		return err
	})
}

func (db *boltDB) DeleteBucket(name string) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(name))
	})
}

func (db *boltDB) Buckets() ([]string, error) {
	buckets := []string{}
	err := db.db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			buckets = append(buckets, string(name))
			return nil
		})
	})
	return buckets, err
}

func (db *boltDB) UseBucket(name string) error {
	nme := []byte(name)
	return db.db.View(func(tx *bolt.Tx) error {
		if b := tx.Bucket(nme); b == nil {
			return ErrBucketNotFound
		}
		db.currBucket = nme
		return nil
	})
}

func (db *boltDB) CurrentBucket() string {
	return string(db.currBucket)
}

func (db *boltDB) Put(key, value string) error {
	if len(key) == 0 {
		return ErrKeyRequired
	}

	if len(value) == 0 {
		return ErrValueRequired
	}

	return db.do(db.db.Update, func(b *bolt.Bucket) error {
		return b.Put([]byte(key), []byte(value))
	})
}

func (db *boltDB) Remove(key string) error {
	if len(key) == 0 {
		return ErrKeyRequired
	}

	return db.do(db.db.Update, func(b *bolt.Bucket) error {
		return b.Delete([]byte(key))
	})
}

func (db *boltDB) Get(key string) (string, error) {
	if len(key) == 0 {
		return "", ErrKeyRequired
	}

	var value string
	err := db.do(db.db.View, func(b *bolt.Bucket) error {
		v := b.Get([]byte(key))
		if len(v) == 0 {
			return ErrKeyNotFound
		}
		value = string(v)
		return nil
	})

	return value, err
}

func (db *boltDB) Keys() ([]string, error) {
	keys := []string{}

	err := db.do(db.db.View, func(b *bolt.Bucket) error {
		return b.ForEach(func(k, _ []byte) error {
			keys = append(keys, string(k))
			return nil
		})
	})

	return keys, err
}

func (db *boltDB) List() (map[string]string, error) {
	kv := map[string]string{}

	err := db.do(db.db.View, func(b *bolt.Bucket) error {
		return b.ForEach(func(k, v []byte) error {
			kv[string(k)] = string(v)
			return nil
		})
	})

	return kv, err
}

func (db *boltDB) do(action func(fn func(*bolt.Tx) error) error, handler func(*bolt.Bucket) error) error {
	return action(func(tx *bolt.Tx) error {
		b := tx.Bucket(db.currBucket)
		if b == nil {
			return ErrBucketNotFound
		}
		return handler(b)
	})
}
