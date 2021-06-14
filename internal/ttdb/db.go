package ttdb

import (
	"encoding/binary"
	"os"
	"os/user"
	"path/filepath"
	"time"
	"path"

	"github.com/alanxoc3/ttrack/internal/seconds"

	bolt "go.etcd.io/bbolt"
)

func getHomeFilePath(filename string) (string, error) {
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir + "/.local/share/ttrack/" + filename, nil
	} else {
		return "", err
	}
}

func open(cacheDir string) (*bolt.DB, error) {
	dbpath := path.Join(cacheDir, "db")

	err := os.MkdirAll(filepath.Dir(dbpath), 0755)
	if err != nil {
		return nil, err
	}

	db, err := bolt.Open(dbpath, 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetSeconds(b *bolt.Bucket, key string) seconds.Seconds {
	return seconds.CreateFromBytes(b.Get([]byte(key)))
}

func SetSeconds(b *bolt.Bucket, key string, secs seconds.Seconds) {
	if secs > seconds.SECONDS_IN_DAY {
		secs = seconds.SECONDS_IN_DAY
	}
	timeout_bytes := make([]byte, 4, 4)
	binary.BigEndian.PutUint32(timeout_bytes[:], uint32(secs))

	if secs == 0 {
		b.Delete([]byte(key))
	} else {
		b.Put([]byte(key), timeout_bytes) // TODO: Error handling.
	}
}

func GetTimestamp(b *bolt.Bucket, key string) time.Time {
	t := time.Time{}
	b_val := b.Get([]byte(key))
	if b_val != nil {
		t.UnmarshalBinary(b_val) // TODO: Error handling.
	}
	return t
}

func SetTimestamp(b *bolt.Bucket, key string, t time.Time) {
	raw, _ := t.MarshalBinary() // TODO: Error handling.
	b.Put([]byte(key), raw)     // TODO: Error handling.
}

func AddTimestampToBucket(b *bolt.Bucket, date_key string, secs seconds.Seconds) {
	old_seconds := GetSeconds(b, date_key)
	SetSeconds(b, date_key, old_seconds+secs)
}

func ViewCmd(cacheDir string, f func(*bolt.Tx) error) {
	db, err := open(cacheDir)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		return f(tx)
	})

	if err != nil {
		panic(err)
	}
}

func UpdateCmd(cacheDir string, f func(*bolt.Tx) error) {
	db, err := open(cacheDir)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		return f(tx)
	})

	if err != nil {
		panic(err)
	}
}

