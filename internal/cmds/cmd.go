package cmds

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
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
        	if b == nil { continue }

            beg_ts, end_ts, timeout := expandGroup(b)
    		_, _, duration, _ := recLogic(s.Now, beg_ts, end_ts, timeout)

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

func getRelWithPanic(basepath, relpath string) string {
	path, err := filepath.Rel(basepath, relpath)
	if err != nil {
		panic(err) // is this possible if data_dir is invalid?
	}
	return path
}

func getListOfGroups(data_dir string, strat walkstrategy) []string {
	groups := []string{}

	filepath.Walk(data_dir, func(path string, info os.FileInfo, e error) error {
		if info == nil {
			return filepath.SkipDir
		}

		group_name := getRelWithPanic(data_dir, path)
		if types.IsValidGroupFile(group_name) || types.IsValidGroupFolder(group_name) {
			group_cleaned_name := types.CreateGroupFromString(group_name).String()
			groups = append(groups, group_cleaned_name)
		}

        if info.IsDir() && group_name != "." && (strat == walk_level || !types.IsValidGroupFolder(group_name)) {
            return filepath.SkipDir
        } else {
    		return nil
        }
	})

	return groups
}

func walkThroughGroups(cache_dir, data_dir string, groupdirs []types.Group, strat walkstrategy) map[types.Group]bool {
	visited_groups := map[types.Group]bool{}

	ttdb.ViewCmd(cache_dir, func(tx *bolt.Tx) error {
		c := tx.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
            cursor_group := types.CreateGroupFromString(string(k))
            for _, group := range groupdirs {
                for _, ancestor := range cursor_group.GetAncestors(group) {
                    visited_groups[ancestor] = true
                    if strat == walk_level && ancestor != group { break }
                }
            }
		}

		return nil
	})

	for _, groupdir := range groupdirs {
    	// Check for the group itself. It could be a folder or a .tt file.
        if !groupdir.IsZero() {
            _, folder_err := os.Stat(filepath.Join(data_dir, groupdir.Filename()));
            _, file_err   := os.Stat(filepath.Join(data_dir, groupdir.String()));

    		if !errors.Is(folder_err, fs.ErrNotExist) || !errors.Is(file_err, fs.ErrNotExist) {
    			if _, exists := visited_groups[groupdir]; !exists {
    				visited_groups[groupdir] = true
    			}
    		}
        }

        // Check the cache. Exact match or begins with .String() + "/".
        // Do I need to split up the groups based on subdirs? Yes I do.
        // How do I keep the order pretty?
        // I probably have to load everything up at startup.

        // Check for all sub groups.
		groupdir_str := groupdir.String()
		groups := getListOfGroups(filepath.Join(data_dir, groupdir_str), strat)
		for _, group := range groups {
			group_with_path := types.CreateGroupFromString(filepath.Join(groupdir_str, group))

			if _, exists := visited_groups[group_with_path]; !exists {
				visited_groups[group_with_path] = true
			}
		}
	}

	return visited_groups
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

func RecFunc(s *State) {
	// Update cache.
	// If new, append to txt file.
	var timestamp_to_write *types.Date
	var seconds_to_write *types.DaySeconds

	group := s.Groups[0]
	timeout_param := s.Duration

	ttdb.UpdateCmd(s.CacheDir, func(tx *bolt.Tx) error {
		// If the bucket doesn't exist and the timeout is zero, do nothing.
		b, err := getOrCreateBucketConditionally(tx, group.String(), timeout_param.IsZero())
		if b == nil || err != nil {
			return err
		}

		old_beg_ts, old_end_ts, old_timeout := expandGroup(b)
		beg_ts, end_ts, duration, finish := recLogic(s.Now, old_beg_ts, old_end_ts, old_timeout)

		if finish && !duration.IsZero() {
			timestamp_to_write = types.CreateDateFromTime(old_beg_ts)
			seconds_to_write = &duration
		}

		if !old_beg_ts.Equal(beg_ts) {
			ttdb.SetTimestamp(b, "beg", beg_ts)
		}
		ttdb.SetTimestamp(b, "end", end_ts)

		if !timeout_param.IsZero() {
			ttdb.SetSeconds(b, "out", timeout_param)
		}

		return nil
	})

	if timestamp_to_write != nil && seconds_to_write != nil {
		ttfile.ModifyTime(filepath.Join(s.DataDir, group.Filename()), *timestamp_to_write, func(ds types.DaySeconds) types.DaySeconds {
			return ds.Add(*seconds_to_write)
		})
	}
}
