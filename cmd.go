package main

import (
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

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

func mvFunc(srcGroup string, dstGroup string) {
	updateCmd(func(tx *bolt.Tx) error {
		srcMap := getDateMap(tx, srcGroup)
		dstMap := getDateMap(tx, dstGroup)

		for k, v := range srcMap {
    			addSecondToMap(dstMap, k, v)
		}

		for k, v := range srcMap {
    			addSecondToMap(dstMap, k, v)
		}

		return nil
	})
}

func setFunc(group string, timestamp time.Time, duration uint32) {
	updateCmd(func(tx *bolt.Tx) error {
		gb, err := tx.CreateBucketIfNotExists([]byte(group))
		if err != nil { return err }

		rb, err := gb.CreateBucketIfNotExists([]byte("rec"))
		if err != nil { return err }

		date_key := formatTimestamp(timestamp)
		setSeconds(rb, date_key, duration)

		return nil
	})
}

func aggFunc(group string) {
	viewCmd(func(tx *bolt.Tx) error {
    	        m := getDateMap(tx, group)
    	        var secs uint32 // This limits the agg output to 136 years. Meh, I won't live that long.
    	        for _, v := range m {
        	        secs += v
    	        }

		fmt.Printf("%s\n", secondsToString(secs))

		return nil
	})
}

func listFunc(group string) {
	viewCmd(func(tx *bolt.Tx) error {
    	        m := getDateMap(tx, group)
    	        for k, v := range m {
			fmt.Printf("%s: %s\n", k, secondsToString(v))
    	        }

		return nil
	})
}

func groupsFunc() {
	viewCmd(func(tx *bolt.Tx) error {
		c := tx.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			fmt.Printf("%s\n", k)
		}

		return nil
	})
}

func delFunc(group string, beg_ts, end_ts time.Time) {
	updateCmd(func(tx *bolt.Tx) error {
		gb := tx.Bucket([]byte(group))
		if gb == nil {
			return nil
		}

		if beg_ts.IsZero() && end_ts.IsZero() {
    			tx.DeleteBucket([]byte(group))
			return nil
		}

		return nil
	})
}

func recFunc(group string, timeout_param uint32) {
	updateCmd(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(group))
		if err != nil {
			return err
		}

		old_beg_ts, old_end_ts, old_timeout, _ := expandGroup(b)
		beg_ts, end_ts, duration, finish := recLogic(time.Now(), old_beg_ts, old_end_ts, old_timeout)

		if finish && duration > 0 {
			rb, err := b.CreateBucketIfNotExists([]byte("rec"))
			if err != nil {
				return err
			}

			addTimestampToBucket(rb, old_beg_ts, duration)
		}

		if !old_beg_ts.Equal(beg_ts) {
			setTimestamp(b, "beg", beg_ts)
		}
		setTimestamp(b, "end", end_ts)

		if timeout_param > 0 {
			setSeconds(b, "out", timeout_param)
		}

		return nil
	})
}
