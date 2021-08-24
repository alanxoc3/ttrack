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
    walk_exist walkstrategy = iota
    walk_level
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

func AggFunc(s *State) {
	groupMap := map[types.Group]bool{}

	walkThroughGroups(s.DataDir, s.Groups, walk_recursive, func(g types.Group) {
		groupMap[g] = true
	})

	date_map := map[types.Date]types.DaySeconds{}
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

func ViewFunc(s *State) {
	/*
		    group := s.Groups[0]
		    beg_date := s.BeginDate.ToDate()
		    end_date := s.EndDate.ToDate()

			dateMap := map[string]types.DaySeconds{}

			ttdb.ViewCmd(s.CacheDir, func(tx *bolt.Tx) error {
				dateMap = getDateMap(tx, group, beg_date.String(), end_date.String())
				return nil
			})

			dates := make([]string, 0, len(dateMap))
			for k := range dateMap {
				dates = append(dates, k)
			}

			sort.Strings(dates)
			for _, d := range dates {
				if v, ok := dateMap[d]; ok {
					fmt.Printf("%s %s\n", d, v.String())
				}
			}
	*/
}

func ListFunc(s *State) {
	groupList := []string{}
	ttdb.ViewCmd(s.CacheDir, func(tx *bolt.Tx) error {
		c := tx.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			groupList = append(groupList, string(k))
		}

		return nil
	})

	sort.Strings(groupList)
	for _, g := range groupList {
		fmt.Printf("%s\n", g)
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
		panic(err) // is this possible if datadir is invalid?
	}
	return path
}

func getListOfGroups(dir string, strat walkstrategy) []string {
	groups := []string{}
	filepath.Walk(dir, func(path string, info os.FileInfo, e error) error {
		if info == nil {
			return filepath.SkipDir
		}

		group_name := getRelWithPanic(dir, path)
		if types.IsValidGroupFile(group_name) || types.IsValidGroupFolder(group_name) {
			group_cleaned_name := types.CreateGroupFromString(group_name).String()
			groups = append(groups, group_cleaned_name)
		}

        if strat == walk_exist {
			return filepath.SkipDir
        } else if strat == walk_level {
            if info.IsDir() && group_name == "." {
                return nil
            } else if info.IsDir() {
    			return filepath.SkipDir
            }
        } else {
            if info.IsDir() && (group_name == "." || types.IsValidGroupFolder(group_name)) {
                return nil
            } else if info.IsDir() {
    			return filepath.SkipDir
            }
        }

		return nil
	})

	return groups
}

func walkThroughGroups(datadir string, groupdirs []types.Group, strat walkstrategy, walkFunc func(types.Group)) {
	visited_groups := map[types.Group]bool{}
	for _, groupdir := range groupdirs {
        if !groupdir.IsZero() {
    		if _, err := os.Stat(filepath.Join(datadir, groupdir.Filename())); !errors.Is(err, fs.ErrNotExist) {
    			if _, exists := visited_groups[groupdir]; !exists {
    				visited_groups[groupdir] = true
    				walkFunc(groupdir)
    			}
    		}

    		if _, err := os.Stat(filepath.Join(datadir, groupdir.String())); !errors.Is(err, fs.ErrNotExist) {
    			if _, exists := visited_groups[groupdir]; !exists {
    				visited_groups[groupdir] = true
    				walkFunc(groupdir)
    			}
    		}
        }

		groupdir_str := groupdir.String()
		groups := getListOfGroups(filepath.Join(datadir, groupdir_str), strat)
		for _, group := range groups {
			group_with_path := types.CreateGroupFromString(filepath.Join(groupdir_str, group))

			if _, exists := visited_groups[group_with_path]; !exists {
				visited_groups[group_with_path] = true
				walkFunc(group_with_path)
			}
		}
	}
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

	walkThroughGroups(s.DataDir, groups, strat, func(path types.Group) {
		fmt.Println(path.String())
	})
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
