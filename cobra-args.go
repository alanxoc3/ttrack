package main

import "github.com/alanxoc3/ttrack/internal/date"

type dateArg date.Date

func (d *dateArg) Set(s string) error {
	v, err := date.CreateFromString(s)
	if err == nil { *d = dateArg(*v) }
	return err
}

func (d *dateArg) Type() string { return "date" }

func (d *dateArg) ToDate() date.Date {
	return (date.Date)(*d)
}

func (d *dateArg) String() string {
	return (*date.Date)(d).String()
}
