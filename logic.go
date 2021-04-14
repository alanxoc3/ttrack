package main

import (
	"fmt"
	"time"
)

func recLogic(now, beg_ts, end_ts time.Time, timeout uint32) (time.Time, time.Time, uint32) {
	time_elapsed := now.Sub(end_ts)
	duration := uint32(0)

	if beg_ts.IsZero() || end_ts.IsZero() {
		fmt.Println("Problem, reseting current time.")
		beg_ts = now
	} else if time_elapsed.Seconds() > float64(timeout) {
		duration = uint32(end_ts.Sub(beg_ts).Seconds()) + timeout
		fmt.Println("Adding time_elapsed: ", duration, "seconds")
		beg_ts = now
	} else {
		fmt.Println("Making current time_elapsed longer:", int64(time_elapsed.Seconds()))
	}

	return beg_ts, now, duration
}
