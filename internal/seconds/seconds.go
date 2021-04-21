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

func CreateFromBytes(b []byte) Seconds {
	if b == nil {
		return 0
	}
	secs := Seconds(binary.BigEndian.Uint32(b))
	if secs > SECONDS_IN_DAY {
		secs = SECONDS_IN_DAY
	}
	return secs
}

func (s *Seconds) String() string {
	dur := time.Duration(*s) * time.Second
	return dur.String()
}
