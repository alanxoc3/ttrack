package main

import (
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

var DEFAULT_TIMEOUT_IN_SECONDS uint32 = 30

func getTimeout(b *bolt.Bucket) uint32 {
	secs := getSeconds(b, "out")
	if secs == 0 {
		return DEFAULT_TIMEOUT_IN_SECONDS
	}
	return secs
}

func setTimeout(b *bolt.Bucket, timeout_in_seconds uint32) {
	setSeconds(b, "out", timeout_in_seconds)
}

func groupsFunc() {
	db, err := opendb()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		c := tx.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			fmt.Printf("%s\n", k)
		}

		return nil
	})
}


func delFunc(group string, beg_ts, end_ts time.Time) {
	db, err := opendb()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
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

	if beg_ts.IsZero() && end_ts.IsZero() {

		fmt.Println("bd: " + beg_ts.String())
		fmt.Println("ed: " + end_ts.String())
	}
}

func recFunc(group string, timeout_param uint32) {
	db, err := opendb()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(group))
		if err != nil {
			return err
		}

		old_beg_ts := getTimestamp(b, "beg")
		beg_ts, end_ts, duration := recLogic(time.Now(), old_beg_ts, getTimestamp(b, "end"), getTimeout(b))

		if duration > 0 {
			rb, err := tx.CreateBucketIfNotExists([]byte("rec"))
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
	if err != nil {
		panic(err)
	}
}
