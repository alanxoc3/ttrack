package date

import "time"
import "fmt"

const DATE_FORMAT_STRING string = "2006-01-02"

type Date struct {
	year uint16
	month uint8
	day uint8
}

func CreateFromString(datestr string) (*Date, error) {
	ts, err := time.Parse(DATE_FORMAT_STRING, datestr)
	if err != nil { return nil, err }
	return CreateFromTime(ts), nil
}

func CreateFromTime(ts time.Time) *Date {
	return &Date{uint16(ts.Year()-1), uint8(ts.Month()-1), uint8(ts.Day()-1)}
}

func dateStrToTimestamp(datestr string) (time.Time, error) {
	return time.Parse(DATE_FORMAT_STRING, datestr)
}

func (d *Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.year+1, d.month+1, d.day+1)
}

func (d *Date) IsZero() bool {
	return d.year == 0 && d.month == 0 && d.day == 0
}

func (d1 *Date) IsLessThan(d2 Date) bool {
	return d1.year < d2.year ||
        d1.year == d2.year && d1.month < d2.month ||
        d1.year == d2.year && d1.month == d2.month && d1.day < d2.day
}

func (d *Date) Set(s string) error {
	v, err := CreateFromString(s)
	if err == nil { *d = Date(*v) }
	return err
}

func (d *Date) Type() string { return "date" }

func (d *Date) ToDate() Date {
	return (Date)(*d)
}
