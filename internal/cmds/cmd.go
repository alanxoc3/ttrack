package cmds

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
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
	walk_files
)

type State struct {
	CacheDir  string
	DataDir   string
	BeginDate types.Date
	EndDate   types.Date
	Recursive bool
	Cached    bool
    Stored    bool
	Daily     bool
	Groups    []types.Group
	Date      types.Date
	Now       time.Time
	Duration  types.DaySeconds
}

// No files are deleted, but files can be added. Files are cleaned too.
func TidyFunc(s *State) string {
    type datesecond struct {
        date types.Date
        seconds types.DaySeconds
    }

    cache := map[types.Group]datesecond{}

    // Clean the cache.
	ttdb.UpdateCmd(s.CacheDir, func(tx *bolt.Tx) error {
		c := tx.Cursor()

		bucketsToDelete := [][]byte{}

		// All items in tx cursor is guaranteed to be a bucket.
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
		    b := tx.Bucket(k)
    		beg_ts, end_ts, timeout := expandGroup(b)
    		if beg_ts.IsZero() || end_ts.IsZero() || timeout.IsZero() {
                bucketsToDelete = append(bucketsToDelete, k)
    		} else {
        		duration, isFinal := calcDuration(s.Now, beg_ts, end_ts, timeout)
        		if isFinal && !duration.IsZero() {
                    bucketsToDelete = append(bucketsToDelete, k)
                    cache[types.CreateGroupFromString(string(k))] = datesecond{*types.CreateDateFromTime(beg_ts), duration}
        		}
    		}
		}

		for _, v := range bucketsToDelete { tx.DeleteBucket(v) }

		return nil
	})

    // Clean existing files.
	groups := walkThroughGroups("", s.DataDir, []types.Group{types.CreateGroupFromString("")}, walk_files, false, true)
	for group := range groups {
		ttfile.ModifyTime(filepath.Join(s.DataDir, group.Filename()), cache[group].date, func(ds types.DaySeconds) types.DaySeconds {
			return ds.Add(cache[group].seconds)
		})
		delete(cache, group)
	}

    // Add new files.
	for group, v := range cache {
		ttfile.ModifyTime(filepath.Join(s.DataDir, group.Filename()), v.date, func(ds types.DaySeconds) types.DaySeconds {
			return ds.Add(v.seconds)
		})
	}

	return ""
}

func SetFunc(s *State) string {
	ttfile.ModifyTime(filepath.Join(s.DataDir, s.Groups[0].Filename()), s.Date, func(ds types.DaySeconds) types.DaySeconds {
		return s.Duration
	})
	return ""
}

func AddFunc(s *State) string {
	ttfile.ModifyTime(filepath.Join(s.DataDir, s.Groups[0].Filename()), s.Date, func(ds types.DaySeconds) types.DaySeconds {
		return ds.Add(s.Duration)
	})
	return ""
}

func SubFunc(s *State) string {
	ttfile.ModifyTime(filepath.Join(s.DataDir, s.Groups[0].Filename()), s.Date, func(ds types.DaySeconds) types.DaySeconds {
		return ds.Sub(s.Duration)
	})
	return ""
}

// TODO: Steps 2 & 3 below could be multi-threaded.
func AggFunc(s *State) string {
    groups := s.Groups
	if len(groups) == 0 {
		groups = []types.Group{types.CreateGroupFromString("")}
	}

	date_map := map[types.Date]types.DaySeconds{}

	// STEP 1: Get the groups recursively.
	groupMap := walkThroughGroups(s.CacheDir, s.DataDir, groups, walk_recursive, s.Cached, s.Stored)

	// STEP 2: Populate from the cache.
    if s.Cached {
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
    }

	// STEP 3: Populate from files.
    if s.Stored {
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
    }

    output := ""
	if s.Daily {
		dates := make(types.DateList, 0, len(date_map))
		for k := range date_map {
			dates = append(dates, k)
		}

		sort.Sort(dates)
		for _, d := range dates {
			if v, ok := date_map[d]; ok && !v.IsZero() {
				output += fmt.Sprintf("%s: %s\n", d.String(), v.String())
			}
		}
	} else {
		totalSeconds := types.MultiDaySeconds{}
		for _, v := range date_map {
			totalSeconds = totalSeconds.AddDaySeconds(v)
		}
		output += fmt.Sprintf("%s\n", totalSeconds.String())
	}

	return output
}

func ResetFunc(s *State) string {
    group := s.Groups[0]
	ttdb.UpdateCmd(s.CacheDir, func(tx *bolt.Tx) error {
		c := tx.Cursor()
		bucketsToDelete := [][]byte{}

		// All items in tx cursor is guaranteed to be a bucket.
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
		    kstr := string(k)
    		if group.IsZero() || kstr == group.String() || strings.HasPrefix(kstr, group.String() + "/") {
        		bucketsToDelete = append(bucketsToDelete, k)
    		}
		}

		for _, v := range bucketsToDelete {
            // Ignoring error code, because there isn't much to do if there is an error.
            tx.DeleteBucket(v)
		}

		return nil
	})

	return ""
}

func LsFunc(s *State) string {
	groups := s.Groups
	if len(groups) == 0 {
		groups = []types.Group{types.CreateGroupFromString("")}
	}

	strat := walk_level
	if s.Recursive {
		strat = walk_recursive
	}

	group_map := walkThroughGroups(s.CacheDir, s.DataDir, groups, strat, s.Cached, s.Stored)
	group_list := make([]string, 0, len(group_map))
	for k := range group_map {
		group_list = append(group_list, k.String())
	}
	sort.Strings(group_list)

    output := ""
	for _, v := range group_list {
		output += fmt.Sprintln(v)
	}
    return output
}

// Updates the cache or text file based on current time.
func RecFunc(s *State) string {
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

		if isFinal && !duration.IsZero() {
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

	return ""
}
