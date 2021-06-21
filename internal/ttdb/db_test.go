package ttdb_test

// Tests create bolt databases in temp directories, then clean those directories when the tests are finished.

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/alanxoc3/ttrack/internal/types"
	"github.com/alanxoc3/ttrack/internal/ttdb"
	"github.com/stretchr/testify/assert"

	"io/ioutil"

	bolt "go.etcd.io/bbolt"
)

func execTest(t *testing.T, testFunc func(string)) {
	dir, err := ioutil.TempDir("", "ttrack-test-")
	defer os.RemoveAll(dir)
	if err != nil {
		t.Fatalf("%s", err)
	}
	testFunc(dir)
}

func TestSeconds(t *testing.T) {
	var testVals = []struct {
		expected string
		input string
	}{
		{"24h0m0s", "25h"} ,
		{"24h0m0s", "24h"} ,
		{"0s"     , "0s"}  ,
		{"1m0s"   , "1m"}  ,
		{"0s"     , "100u"},
		{"0s"     , "100o"},
	}

	for _, v := range testVals {
		t.Run(fmt.Sprintf("test-%s", v.input), func(t *testing.T) {
        	execTest(t, func(dir string) {
        		var secs types.Seconds
        		ttdb.UpdateCmd(dir, func(b *bolt.Tx) error {
        			bucket, err := b.CreateBucket([]byte("group"))
        			ttdb.SetSeconds(bucket, "key", types.CreateSecondsFromString(v.input))
        			return err
        		})

        		ttdb.ViewCmd(dir, func(b *bolt.Tx) error {
        			bucket := b.Bucket([]byte("group"))
        			secs = ttdb.GetSeconds(bucket, "key")
        			return nil
        		})

        		assert.Equal(t, v.expected, secs.String())
        	})
		})
	}
}

func TestTimestamp(t *testing.T) {
	var testVals = []struct {
		expected string
	}{
		{"0001-01-01T00:00:00Z"},
		{"2006-01-02T15:04:05Z"},
		{"2021-12-31T23:59:59Z"},
		{"1999-06-20T15:32:43Z"},
	}
	for _, v := range testVals {
		t.Run(fmt.Sprintf("test-%s", v.expected), func(t *testing.T) {
			expectedTime, err := time.Parse(time.RFC3339, v.expected)
			if err != nil {
				t.Fatalf("%s", err)
			}

			execTest(t, func(dir string) {
				ttdb.UpdateCmd(dir, func(b *bolt.Tx) error {
					bucket, err := b.CreateBucket([]byte("group"))
					ttdb.SetTimestamp(bucket, "key", expectedTime)
					return err
				})

				var actualTime time.Time
				ttdb.ViewCmd(dir, func(b *bolt.Tx) error {
					bucket := b.Bucket([]byte("group"))
					actualTime = ttdb.GetTimestamp(bucket, "key")
					return nil
				})

				assert.Equal(t, expectedTime, actualTime)
			})
		})
	}
}
