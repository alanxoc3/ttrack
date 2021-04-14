package main

import (
	"os/user"
	"time"
)

type date time.Time

func newDate(val time.Time, p *time.Time) *date {
	*p = val
	return (*date)(p)
}

func (d *date) Set(s string) error {
	v, err := time.Parse(DATE_FORMAT_STRING, s)
	*d = date(v)
	return err
}

func (d *date) Type() string { return "date" }

func (d *date) String() string {
	if (*time.Time)(d).IsZero() {
		return ""
	} else {
		return (*time.Time)(d).Format(DATE_FORMAT_STRING)
	}
}

func getHomeFilePath(filename string) (string, error) {
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir + "/.config/ttrack/" + filename, nil
	} else {
		return "", err
	}
}
