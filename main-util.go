package main

import (
	"os/user"
	"time"
	"strings"

	bolt "go.etcd.io/bbolt"
)

type date time.Time

func dateStrToTimestamp(datestr string) (time.Time, error) {
	return time.Parse(DATE_FORMAT_STRING, datestr)
}

func (d *date) Set(s string) error {
	v, err := dateStrToTimestamp(s)
	*d = date(v)
	return err
}

func (d *date) Type() string { return "date" }

func (d *date) String() string {
	return formatTimestamp(*(*time.Time)(d))
}

func getHomeFilePath(filename string) (string, error) {
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir + "/.local/share/ttrack/" + filename, nil
	} else {
		return "", err
	}
}

func viewCmd(f func(*bolt.Tx) error) {
	db, err := opendb()
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

func updateCmd(f func(*bolt.Tx) error) {
	db, err := opendb()
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

func clean_group(group string) string {
	fields := strings.FieldsFunc(group, func(c rune) bool { return c == '/' })

	newFields := []string{}
	for _, v := range fields {
    		for i, r := range v {
			if r != '.' {
				newFields = append(newFields, v[i:])
				break
			}
    		}
	}

	group = strings.Join(newFields, "/")
	return group
}
