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

func getDateMap(tx *bolt.Tx, group string) map[string]uint32 {
	m := map[string]uint32{}

	gb := tx.Bucket([]byte(group))
	if gb == nil { return m }

	beg_ts, _, duration, _ := recLogic(time.Now(), getTimestamp(gb, "beg"), getTimestamp(gb, "end"), getTimeout(gb))
	addSecondToMap(m, formatTimestamp(beg_ts), duration)

	rb := gb.Bucket([]byte("rec"))
	if rb == nil { return m }

	c := rb.Cursor()

	for k, v := c.First(); k != nil; k, v = c.Next() {
		addSecondToMap(m, string(k), bytesToSeconds(v))
	}

	return m
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

		old_beg_ts := getTimestamp(b, "beg")
		beg_ts, end_ts, duration, finish := recLogic(time.Now(), old_beg_ts, getTimestamp(b, "end"), getTimeout(b))

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
			setTimeout(b, timeout_param)
		}

		fmt.Println("beg_ts:", beg_ts)
		fmt.Println("end_ts:", end_ts)

		return nil
	})
}
