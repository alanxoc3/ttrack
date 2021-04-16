package main

import (
	"time"
	bolt "go.etcd.io/bbolt"
)

type bucketInterface interface {
	Bucket([]byte)*bolt.Bucket
	CreateBucket([]byte)(*bolt.Bucket, error)
}

func getOrCreateBucketConditionally(parent bucketInterface, key string, nilCondition bool) (*bolt.Bucket, error) {
	b := parent.Bucket([]byte(key))
	if b == nil && nilCondition {
		return nil, nil
	} else if b == nil {
		var err error
		b, err = parent.CreateBucket([]byte(key))
		if err != nil { return nil, err }
	}
	return b, nil
}

func recLogic(now, beg_ts, end_ts time.Time, timeout uint32) (time.Time, time.Time, uint32, bool) {
	time_elapsed := now.Sub(end_ts)
	duration := uint32(0)
	finish := false

	if beg_ts.IsZero() || end_ts.IsZero() {
		beg_ts = now
	} else if time_elapsed.Seconds() > float64(timeout) {
		duration = uint32(end_ts.Sub(beg_ts).Seconds()) + timeout
		finish = true
		beg_ts = now
	} else {
		duration = uint32(now.Sub(beg_ts).Seconds())
	}

	return beg_ts, now, duration, finish
}

func addSecondToMap(m map[string]uint32, key string, num uint32) {
	if num == 0 { return }
	var base_val uint32
	if v, ok := m[key]; ok { base_val = v }
	m[key] = base_val + num
	if m[key] > SECONDS_IN_DAY {
		m[key] = SECONDS_IN_DAY
	}
}

func getGroupBucket(tx *bolt.Tx, group string) *bolt.Bucket {
	return tx.Bucket([]byte(group))
}

func expandGroup(b *bolt.Bucket) (time.Time, time.Time, uint32, *bolt.Bucket) {
	return getTimestamp(b, "beg"), getTimestamp(b, "end"), getSeconds(b, "out"), b.Bucket([]byte("rec"))
}

func getGroupRecBucket(tx *bolt.Tx, group string) *bolt.Bucket {
	gb := tx.Bucket([]byte(group))
	if gb == nil { return nil }

	rb := gb.Bucket([]byte("rec"))
	return rb
}

func getDateMap(tx *bolt.Tx, group string) map[string]uint32 {
	m := map[string]uint32{}

	gb := getGroupBucket(tx, group)
	if gb == nil { return m }

	beg_ts, end_ts, timeout, rec := expandGroup(gb)
	_, _, duration, _ := recLogic(time.Now(), beg_ts, end_ts, timeout)
	addSecondToMap(m, formatTimestamp(beg_ts), duration)

	if rec == nil { return m }

	c := rec.Cursor()

	for k, v := c.First(); k != nil; k, v = c.Next() {
		addSecondToMap(m, string(k), bytesToSeconds(v))
	}

	return m
}

