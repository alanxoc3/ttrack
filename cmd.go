package main

import (
	"fmt"
	"sort"
	"time"

	bolt "go.etcd.io/bbolt"
)

func cpFunc(srcGroup string, dstGroup string) {
	updateCmd(func(tx *bolt.Tx) error {
		m := getDateMap(tx, srcGroup)
		if len(m) == 0 {
			return nil
		}

		dstBucket, err := tx.CreateBucketIfNotExists([]byte(dstGroup))
		if err != nil {
			return err
		}

		rec, err := dstBucket.CreateBucketIfNotExists([]byte("rec"))
		if err != nil {
			return err
		}

		for k, v := range m {
			addTimestampToBucket(rec, k, v)
		}
		return nil
	})
}

func setFunc(group string, timestamp time.Time, duration uint32) {
	updateCmd(func(tx *bolt.Tx) error {
		gb, err := getOrCreateBucketConditionally(tx, group, duration == 0)
		if gb == nil || err != nil {
			return err
		}

		rb, err := getOrCreateBucketConditionally(gb, "rec", duration == 0)
		if rb == nil || err != nil {
			return err
		}

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

func listFunc(group string, beg_ts, end_ts time.Time) {
	dateMap := map[string]uint32{}

	viewCmd(func(tx *bolt.Tx) error {
		dateMap = getDateMap(tx, group)
		return nil
	})

	dates := make([]string, 0, len(dateMap))
	for k := range dateMap {
		dates = append(dates, k)
	}

	sort.Strings(dates)
	for _, d := range dates {
		fmt.Printf("%s %s\n", d, secondsToString(dateMap[d]))
	}
}

func groupsFunc() {
	groupList := []string{}
	viewCmd(func(tx *bolt.Tx) error {
		c := tx.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			groupList = append(groupList, string(k))
		}

		return nil
	})

	sort.Strings(groupList)
	for _, g := range groupList {
		fmt.Printf("%s\n", g)
	}

}

func delFunc(group string) {
	updateCmd(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(group))
		if b == nil {
			return nil
		}

		tx.DeleteBucket([]byte(group))
		return nil
	})
}

func recFunc(group string, timeout_param uint32) {
	updateCmd(func(tx *bolt.Tx) error {
		b, err := getOrCreateBucketConditionally(tx, group, timeout_param == 0)
		if b == nil || err != nil {
			return err
		}

		old_beg_ts, old_end_ts, old_timeout, _ := expandGroup(b)
		beg_ts, end_ts, duration, finish := recLogic(time.Now(), old_beg_ts, old_end_ts, old_timeout)

		if finish && duration > 0 {
			rb, err := b.CreateBucketIfNotExists([]byte("rec"))
			if err != nil {
				return err
			}

			addTimestampToBucket(rb, formatTimestamp(old_beg_ts), duration)
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
