package seconds

import (
	"encoding/binary"
	"time"
)

// This limits the agg output to 136 years. Meh, I won't live that long to care.
type Seconds uint32

const SECONDS_IN_DAY Seconds = 86400

func CreateFromDuration(d time.Duration) Seconds {
	return Seconds(d.Milliseconds() / 1000)
}

// Used with files.
func CreateFromString(s string) Seconds {
    dur, _ := time.ParseDuration(s) // TODO: Should something be done with an error?
    return CreateFromDuration(dur)
}

// Used with the cache.
func CreateFromBytes(b []byte) Seconds {
	if b == nil {
		return 0
	}
	return Seconds(binary.BigEndian.Uint32(b))
}

// Used with files and cache.
func (s *Seconds) CapAtOneDay() Seconds {
	if *s > SECONDS_IN_DAY {
		return SECONDS_IN_DAY
	}
	return *s
}

// Used for writing to files.
func (s *Seconds) String() string {
	dur := time.Duration(*s) * time.Second
	return dur.String()
}
