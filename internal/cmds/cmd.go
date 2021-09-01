package cmds

import (
	"fmt"
	"path/filepath"
	"sort"
	"time"

	"github.com/alanxoc3/ttrack/internal/ttdb"
	"github.com/alanxoc3/ttrack/internal/ttfile"
	"github.com/alanxoc3/ttrack/internal/types"

	bolt "go.etcd.io/bbolt"
)

type walkstrategy uint8

const (
	walk_level walkstrategy = iota
	walk_recursive
)

type State struct {
	CacheDir  string
	DataDir   string
	BeginDate types.Date
	EndDate   types.Date
	Recursive bool
	Daily     bool
	Groups    []types.Group
	Date      types.Date
	Now       time.Time
	Duration  types.DaySeconds
}

func CpFunc(s *State) {
	/*
		    srcGroup := s.Groups[0]
		    dstGroup := s.Groups[0]
		    beg_date := s.BeginDate.ToDate()
		    end_date := s.EndDate.ToDate()

			ttdb.UpdateCmd(s.CacheDir, func(tx *bolt.Tx) error {
				m := getDateMap(tx, srcGroup, beg_date.String(), end_date.String())
				if len(m) == 0 {
					return nil
				}

				dstBucket, err := tx.CreateBucketIfNotExists([]byte(dstGroup))
				if err != nil {
					return err
				}

				rec, err := dstBucket.CreateBucketIfNotExists([]byte("rec"))
				if err != nil {
					return err
				}

				for k, v := range m {
					ttdb.AddTimestampToBucket(rec, k, v)
				}
				return nil
			})
	*/
}

func TidyFunc(s *State) {
	// What does tidy do?
	// Anything in the cache that has exceeded its value goes to a file and that cache entry is deleted.
	// All files are checked, sorted, and merging duplicates.
}

func SetFunc(s *State) {
	ttfile.ModifyTime(filepath.Join(s.DataDir, s.Groups[0].Filename()), s.Date, func(ds types.DaySeconds) types.DaySeconds {
		return s.Duration
	})
}

func AddFunc(s *State) {
	ttfile.ModifyTime(filepath.Join(s.DataDir, s.Groups[0].Filename()), s.Date, func(ds types.DaySeconds) types.DaySeconds {
		return ds.Add(s.Duration)
	})
}

func SubFunc(s *State) {
	ttfile.ModifyTime(filepath.Join(s.DataDir, s.Groups[0].Filename()), s.Date, func(ds types.DaySeconds) types.DaySeconds {
		return ds.Sub(s.Duration)
	})
}

// TODO: Steps 2 & 3 below could be multi-threaded.
func AggFunc(s *State) {
	date_map := map[types.Date]types.DaySeconds{}

	// STEP 1: Get the groups recursively.
	groupMap := walkThroughGroups(s.CacheDir, s.DataDir, s.Groups, walk_recursive)

	// STEP 2: Populate from the cache.
	ttdb.ViewCmd(s.CacheDir, func(tx *bolt.Tx) error {
		for group := range groupMap {
			b := tx.Bucket([]byte(group.String()))
			if b == nil {
				continue
			}

			beg_ts, end_ts, timeout := expandGroup(b)
			duration, _ := calcDuration(s.Now, beg_ts, end_ts, timeout)

			date_key := *types.CreateDateFromTime(beg_ts)

			if types.IsDateBetween(s.BeginDate, date_key, s.EndDate) {
				if date_map_val, exists := date_map[date_key]; exists {
					date_map[date_key] = date_map_val.Add(duration)
				} else {
					date_map[date_key] = duration
				}
			}
		}

		return nil
	})

	// STEP 3: Populate from files.
	for group := range groupMap {
		local_date_map := ttfile.GetDateSeconds(filepath.Join(s.DataDir, group.Filename()))

		for k, v := range local_date_map {
			if types.IsDateBetween(s.BeginDate, k, s.EndDate) {
				if date_map_val, exists := date_map[k]; exists {
					date_map[k] = date_map_val.Add(v)
				} else {
					date_map[k] = v
				}
			}
		}
	}

	if s.Daily {
		dates := make(types.DateList, 0, len(date_map))
		for k := range date_map {
			dates = append(dates, k)
		}

		sort.Sort(dates)
		for _, d := range dates {
			if v, ok := date_map[d]; ok && !v.IsZero() {
				fmt.Printf("%s: %s\n", d.String(), v.String())
			}
		}
	} else {
		totalSeconds := types.MultiDaySeconds{}
		for _, v := range date_map {
			totalSeconds = totalSeconds.AddDaySeconds(v)
		}
		fmt.Printf("%s\n", totalSeconds.String())
	}
}

func DelFunc(s *State) {
	group := s.Groups[0]
	ttdb.UpdateCmd(s.CacheDir, func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(group.String()))
		if b == nil {
			return nil
		}

		tx.DeleteBucket([]byte(group.String()))
		return nil
	})
}

func LsFunc(s *State) {
	groups := s.Groups
	if len(groups) == 0 {
		groups = []types.Group{types.CreateGroupFromString("")}
	}

	strat := walk_level
	if s.Recursive {
		strat = walk_recursive
	}

	group_map := walkThroughGroups(s.CacheDir, s.DataDir, groups, strat)
	group_list := make([]string, 0, len(group_map))
	for k := range group_map {
		group_list = append(group_list, k.String())
	}
	sort.Strings(group_list)

	for _, v := range group_list {
		fmt.Println(v)
	}

}

// Updates the cache or text file based on current time.
func RecFunc(s *State) {
	var timestamp_to_write *types.Date
	var seconds_to_write *types.DaySeconds

	group := s.Groups[0]
	timeout_param := s.Duration

	ttdb.UpdateCmd(s.CacheDir, func(tx *bolt.Tx) error {
        // STEP 1: Create bucket if timeout is non-zero/valid.
    	b := tx.Bucket([]byte(group.String()))
    	if b == nil {
            if timeout_param.IsZero() { return nil }
    		var err error
    		b, err = tx.CreateBucket([]byte(group.String()))
    		if err != nil { return err }
    	}

        // STEP 2: Calc duration.
		beg_ts, end_ts, timeout := expandGroup(b)
		duration, isFinal := calcDuration(s.Now, beg_ts, end_ts, timeout)

		if !duration.IsZero() {
			timestamp_to_write = types.CreateDateFromTime(beg_ts)
			seconds_to_write = &duration
		}

        // STEP 3: Set beg, end, and out.
		ttdb.SetTimestamp(b, "end", s.Now)
		ttdb.SetSeconds(b, "out", timeout_param)

		if isFinal || beg_ts.IsZero() || end_ts.IsZero() || timeout.IsZero() {
			ttdb.SetTimestamp(b, "beg", s.Now)
		}

		return nil
	})

	if timestamp_to_write != nil && seconds_to_write != nil {
		ttfile.ModifyTime(filepath.Join(s.DataDir, group.Filename()), *timestamp_to_write, func(ds types.DaySeconds) types.DaySeconds {
			return ds.Add(*seconds_to_write)
		})
	}
}
