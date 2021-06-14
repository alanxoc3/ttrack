package ttdb_test

// Tests create bolt databases in temp directories, then clean those directories when the tests are finished.

import (
	"os"
	"testing"
	"time"

	"github.com/alanxoc3/ttrack/internal/seconds"
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

func TestSetGetSeconds(t *testing.T) {
	execTest(t, func(dir string) {
        secs := seconds.CreateFromDuration(time.Duration(0))
		ttdb.UpdateCmd(dir, func(b *bolt.Tx) error {
            bucket, err := b.CreateBucket([]byte("group"))
			ttdb.SetSeconds(bucket, "key", seconds.CreateFromDuration(time.Hour))
			return err
		})

		ttdb.ViewCmd(dir, func(b *bolt.Tx) error {
            bucket := b.Bucket([]byte("group"))
			secs = ttdb.GetSeconds(bucket, "key")
			return nil
		})

		assert.Equal(t, seconds.CreateFromDuration(time.Hour), secs)
	})
}

// func GetSeconds(b *bolt.Bucket, key string) seconds.Seconds
// func SetSeconds(b *bolt.Bucket, key string, secs seconds.Seconds) {
// func GetTimestamp(b *bolt.Bucket, key string) time.Time {
// func SetTimestamp(b *bolt.Bucket, key string, t time.Time) {
// func AddTimestampToBucket(b *bolt.Bucket, date_key string, secs seconds.Seconds) {
// func ViewCmd(cacheDir string, f func(*bolt.Tx) error) {
// func UpdateCmd(cacheDir string, f func(*bolt.Tx) error) {
