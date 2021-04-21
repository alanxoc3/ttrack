package main

import (
	"os/user"
	"strings"

	bolt "go.etcd.io/bbolt"
)

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
