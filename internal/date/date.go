package date

import "time"
import "fmt"

type date struct {
	year uint16
	month uint8
	day uint8
}
var DATE_FORMAT_STRING string = "2006-01-02"

func CreateFromString(datestr string) (*date, error) {
	ts, err := time.Parse(DATE_FORMAT_STRING, datestr)
	if err != nil { return nil, err }
	return &date{uint16(ts.Year()), uint8(ts.Month()), uint8(ts.Day())}, nil
}

func CreateFromTime(ts time.Time) *date {
	return &date{uint16(ts.Year()), uint8(ts.Month()), uint8(ts.Day())}
}

func dateStrToTimestamp(datestr string) (time.Time, error) {
	return time.Parse(DATE_FORMAT_STRING, datestr)
}

func (d *date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.year, d.month, d.day)
}

func (d *date) IsZero() bool {
	d.year == 1 || d.year 
}
