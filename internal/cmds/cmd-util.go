package cmds

import (
	"errors"
	"io/fs"
	"os"
	"strings"
	"time"

	"path/filepath"

	"github.com/alanxoc3/ttrack/internal/ttdb"
	"github.com/alanxoc3/ttrack/internal/types"

	bolt "go.etcd.io/bbolt"
)

type bucketInterface interface {
	Bucket([]byte) *bolt.Bucket
	CreateBucket([]byte) (*bolt.Bucket, error)
}

// duration -> the duration that should be saved.
// final -> is the duration final/should beg_ts restart?
func calcDuration(now, beg_ts, end_ts time.Time, timeout types.DaySeconds) (types.DaySeconds, bool) {
    if beg_ts.IsZero() || end_ts.IsZero() || timeout.IsZero() {
        return types.DaySeconds{}, false
    } else if now.Sub(end_ts).Seconds() > float64(timeout.GetAsUint32()) {
		return types.CreateSecondsFromUint32(uint32(end_ts.Sub(beg_ts).Seconds())).Add(timeout), true
	} else {
    	return types.CreateSecondsFromUint32(uint32(now.Sub(beg_ts).Seconds())), false
	}
}

func addSecondToMap(m map[string]types.DaySeconds, key string, num types.DaySeconds) {
	if num.IsZero() {
		return
	}
	var base_val types.DaySeconds
	if v, ok := m[key]; ok {
		base_val = v
	}
	m[key] = base_val.Add(num)
}

func getGroupBucket(tx *bolt.Tx, group string) *bolt.Bucket {
	return tx.Bucket([]byte(group))
}

func expandGroup(b *bolt.Bucket) (time.Time, time.Time, types.DaySeconds) {
	return ttdb.GetTimestamp(b, "beg"), ttdb.GetTimestamp(b, "end"), ttdb.GetSeconds(b, "out")
}

func getGroupRecBucket(tx *bolt.Tx, group string) *bolt.Bucket {
	gb := tx.Bucket([]byte(group))
	if gb == nil {
		return nil
	}

	rb := gb.Bucket([]byte("rec"))
	return rb
}

func is_date_str_in_range(date, beg_date, end_date string) bool {
	return (beg_date == "" || strings.Compare(beg_date, date) <= 0) &&
		(end_date == "" || strings.Compare(end_date, date) >= 0)
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
					if strat == walk_level && ancestor != group {
						break
					}
				}
			}
		}

		return nil
	})

	for _, groupdir := range groupdirs {
		// Check for the group itself. It could be a folder or a .tt file.
		if !groupdir.IsZero() {
			_, folder_err := os.Stat(filepath.Join(data_dir, groupdir.Filename()))
			_, file_err := os.Stat(filepath.Join(data_dir, groupdir.String()))

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

/*
func getDateMap(tx *bolt.Tx, group, beg_bounds, end_bounds string) map[string]types.DaySeconds {
	m := map[string]types.DaySeconds{}

	gb := getGroupBucket(tx, group)
	if gb == nil { return m }

	beg_ts, end_ts, timeout, := expandGroup(gb)
	_, _, duration, _ := recLogic(time.Now(), beg_ts, end_ts, timeout)
	beg_ts_str := formatTimestamp(beg_ts)
	if is_date_str_in_range(beg_ts_str, beg_bounds, end_bounds) {
		addSecondToMap(m, formatTimestamp(beg_ts), duration)
	}

	if rec == nil { return m }

	c := rec.Cursor()

	for k, v := c.First(); k != nil; k, v = c.Next() {
    		dateStr := string(k)
    		if is_date_str_in_range(dateStr, beg_bounds, end_bounds) {
			addSecondToMap(m, string(k), types.CreateFromBytes(v))
    		}
	}

	return m
}
*/
