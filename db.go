package main

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"time"

	bolt "go.etcd.io/bbolt"
)

var DATE_FORMAT_STRING string = "2006-01-02"
var SECONDS_IN_DAY uint32 = 86400

func opendb() (*bolt.DB, error) {
	dbpath, err := getHomeFilePath("db")
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(filepath.Dir(dbpath), 0755)
	if err != nil {
		return nil, err
	}

	db, err := bolt.Open(dbpath, 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getSeconds(b *bolt.Bucket, key string) uint32 {
	var secs uint32 = 0
	b_val := b.Get([]byte(key))
	if b_val != nil {
		secs = binary.BigEndian.Uint32(b_val)
	}

	if secs > SECONDS_IN_DAY {
		secs = SECONDS_IN_DAY
	}

	return secs
}

func setSeconds(b *bolt.Bucket, key string, secs uint32) {
	timeout_bytes := make([]byte, 4, 4)
	binary.BigEndian.PutUint32(timeout_bytes[:], secs)
	b.Put([]byte(key), timeout_bytes) // TODO: Error handling.
}

func getTimestamp(b *bolt.Bucket, key string) time.Time {
	t := time.Time{}
	b_val := b.Get([]byte(key))
	if b_val != nil {
		t.UnmarshalBinary(b_val) // TODO: Error handling.
	}
	return t
}

func setTimestamp(b *bolt.Bucket, key string, t time.Time) {
	raw, _ := t.MarshalBinary() // TODO: Error handling.
	b.Put([]byte(key), raw)     // TODO: Error handling.
}

func addTimestampToBucket(b *bolt.Bucket, beg_ts time.Time, seconds uint32) {
	date_key := beg_ts.Format(DATE_FORMAT_STRING)
	old_seconds := getSeconds(b, date_key)
	setSeconds(b, date_key, old_seconds+seconds)
}
