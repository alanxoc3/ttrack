package types

import (
	"encoding/binary"
	"time"
)

// Limit for agg output is like 5 billion centuries. So yeah... a long long time.
type MultiDaySeconds struct {
	val uint64
}

func (mds MultiDaySeconds) AddDaySeconds(ds DaySeconds) MultiDaySeconds {
    return MultiDaySeconds{val: mds.val + uint64(ds.val)}
}

func (s *MultiDaySeconds) String() string {
	dur := time.Duration(s.val) * time.Second
	return dur.String()
}

// DaySeconds Below

type DaySeconds struct {
	val uint32
}

const seconds_in_day uint32 = 86400

func CreateSecondsFromDuration(d time.Duration) DaySeconds {
	return DaySeconds{val: uint32(d.Milliseconds()/1000)}.capAtOneDay()
}

func CreateSecondsFromString(s string) DaySeconds {
	dur, _ := time.ParseDuration(s) // TODO: Should something be done with an error?
	return CreateSecondsFromDuration(dur)
}

func CreateSecondsFromBytes(b []byte) DaySeconds {
	if b == nil {
		return DaySeconds{val: 0}
	}
	return DaySeconds{val: binary.BigEndian.Uint32(b)}.capAtOneDay()
}

func CreateSecondsFromUint32(num uint32) DaySeconds {
	return DaySeconds{val: num}.capAtOneDay()
}

func (s DaySeconds) GetAsUint32() uint32 {
    return s.val
}

func (s DaySeconds) Add(s2 DaySeconds) DaySeconds {
    return DaySeconds{val: s.val + s2.val}.capAtOneDay()
}

func (s DaySeconds) Sub(s2 DaySeconds) DaySeconds {
    // Capping at one day isn't needed for subtraction.
    if s.val <= s2.val {
        return DaySeconds{0}
    } else {
        return DaySeconds{val: s.val - s2.val}
    }
}

// Used for writing to files.
func (s *DaySeconds) String() string {
	dur := time.Duration(s.val) * time.Second
	return dur.String()
}

func (s DaySeconds) IsZero() bool {
    return s.val == 0
}

func (s DaySeconds) capAtOneDay() DaySeconds {
	if s.val > seconds_in_day {
		return DaySeconds{val: seconds_in_day}
	}
	return s
}

