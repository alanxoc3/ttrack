package main

import (
	"time"
)

func recLogic(now, beg_ts, end_ts time.Time, timeout uint32) (time.Time, time.Time, uint32, bool) {
	time_elapsed := now.Sub(end_ts)
	duration := uint32(0)
	finish := false

	if beg_ts.IsZero() || end_ts.IsZero() {
		beg_ts = now
	} else if time_elapsed.Seconds() > float64(timeout) {
		duration = uint32(end_ts.Sub(beg_ts).Seconds()) + timeout
		finish = true
		beg_ts = now
	} else {
		duration = uint32(now.Sub(beg_ts).Seconds())
	}

	return beg_ts, now, duration, finish
}
