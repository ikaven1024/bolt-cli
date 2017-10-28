package db

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDB_Bucket(t *testing.T) {
	path := filepath.Join(os.TempDir(), "bolt_test_"+time.Now().Format("20060102150405"))
	db, err := New(path)
	require.NoError(t, err, "new bolt db")
	defer os.Remove(path)
	defer db.Close()

	buckets, err := db.Buckets()
	assert.NoError(t, err, "buckets")
	assert.Len(t, buckets, 0, "length of bucket before create")

	assert.NoError(t, db.CreateBucket("test-bucket-1"), "create bucket 1")
	assert.NoError(t, db.CreateBucket("test-bucket-2"), "create bucket 2")
	assert.NoError(t, db.CreateBucket("test-bucket-3"), "create bucket 3")
	assert.Equal(t, ErrBucketExists, db.CreateBucket("test-bucket-1"), "create bucket exist")

	assert.NoError(t, db.DeleteBucket("test-bucket-1"), "delete bucket 1")
	assert.Equal(t, ErrBucketNotFound, db.DeleteBucket("non-exist-bucket"), "delete bucket non-exist")

	buckets, err = db.Buckets()
	assert.NoError(t, err, "buckets")
	assert.Len(t, buckets, 2, "length of bucket after create")

	assert.NoError(t, db.UseBucket("test-bucket-2"), "use bucket test-bucket-1")
	assert.Equal(t, ErrBucketNotFound, db.UseBucket("non-exist-bucket"), "use bucket non-exist")

	assert.Equal(t, "test-bucket-2", db.CurrentBucket(), "current bucket")
}

func TestDB_KeyValue(t *testing.T) {
	db, cancel := createTestDB()
	defer cancel()

	t.Log()
	assert.NoError(t, db.Put("key1", "value1"), "put key1")
	assert.NoError(t, db.Put("key2", "value2"), "put key2")
	assert.NoError(t, db.Put("key3", "value3"), "put key3")

	assert.NoError(t, db.Put("key1", "value1_new"), "update key1")

	assert.NoError(t, db.Remove("key2"), "delete key2")
	assert.NoError(t, db.Remove("key-non-exist"), "delete non-exist key")

	v1, err := db.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1_new", v1, "value of key1 aafer update")

	_, err = db.Get("key-non-exist")
	assert.Equal(t, ErrKeyNotFound, err, "get non-exist key")

	keys, err := db.Keys()
	assert.NoError(t, err)
	assert.Len(t, keys, 2, "length of keys")

	kv, err := db.List()
	assert.NoError(t, err)
	assert.Len(t, kv, 2, "length of keys")
}

func createTestDB() (*boltDB, func()) {
	path := filepath.Join(os.TempDir(), "bolt_test_"+time.Now().Format("20060102150405"))
	db, _ := New(path)
	db.CreateBucket("test")
	db.UseBucket("test")

	return db, func() {
		db.Close()
		os.Remove(path)
	}
}
