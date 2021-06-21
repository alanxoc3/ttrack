package cmds

import "github.com/alanxoc3/ttrack/internal/types"
import "github.com/alanxoc3/ttrack/internal/ttdb"
import (
	"strings"
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

func recLogic(now, beg_ts, end_ts time.Time, timeout types.Seconds) (time.Time, time.Time, types.Seconds, bool) {
	time_elapsed := now.Sub(end_ts)
	duration := types.Seconds(0)
	finish := false

	if beg_ts.IsZero() || end_ts.IsZero() {
		beg_ts = now
	} else if time_elapsed.Seconds() > float64(timeout) {
		duration = types.Seconds(end_ts.Sub(beg_ts).Seconds()) + timeout
		finish = true
		beg_ts = now
	} else {
		duration = types.Seconds(now.Sub(beg_ts).Seconds())
	}

	return beg_ts, now, duration, finish
}

func addSecondToMap(m map[string]types.Seconds, key string, num types.Seconds) {
	if num == 0 { return }
	var base_val types.Seconds
	if v, ok := m[key]; ok { base_val = v }
	m[key] = base_val + num
	if m[key] > types.SECONDS_IN_DAY {
		m[key] = types.SECONDS_IN_DAY
	}
}

func getGroupBucket(tx *bolt.Tx, group string) *bolt.Bucket {
	return tx.Bucket([]byte(group))
}

func expandGroup(b *bolt.Bucket) (time.Time, time.Time, types.Seconds) {
	return ttdb.GetTimestamp(b, "beg"), ttdb.GetTimestamp(b, "end"), ttdb.GetSeconds(b, "out")
}

func getGroupRecBucket(tx *bolt.Tx, group string) *bolt.Bucket {
	gb := tx.Bucket([]byte(group))
	if gb == nil { return nil }

	rb := gb.Bucket([]byte("rec"))
	return rb
}

func is_date_str_in_range(date, beg_date, end_date string) bool {
	return (beg_date == "" || strings.Compare(beg_date, date) <= 0) &&
		(end_date == "" || strings.Compare(end_date, date) >= 0)
}

/*
func getDateMap(tx *bolt.Tx, group, beg_bounds, end_bounds string) map[string]types.Seconds {
	m := map[string]types.Seconds{}

	gb := getGroupBucket(tx, group)
	if gb == nil { return m }

	beg_ts, end_ts, timeout, := expandGroup(gb)
	_, _, duration, _ := recLogic(time.Now(), beg_ts, end_ts, timeout)
	beg_ts_str := formatTimestamp(beg_ts)
	if is_date_str_in_range(beg_ts_str, beg_bounds, end_bounds) {
		addSecondToMap(m, formatTimestamp(beg_ts), duration)
	}

	if rec == nil { return m }

	c := rec.Cursor()

	for k, v := c.First(); k != nil; k, v = c.Next() {
    		dateStr := string(k)
    		if is_date_str_in_range(dateStr, beg_bounds, end_bounds) {
			addSecondToMap(m, string(k), types.CreateFromBytes(v))
    		}
	}

	return m
}
*/
