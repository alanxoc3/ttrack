package cmds

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/alanxoc3/ttrack/internal/ttdb"
	"github.com/alanxoc3/ttrack/internal/ttfile"
	"github.com/alanxoc3/ttrack/internal/types"

	bolt "go.etcd.io/bbolt"
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
	group := s.Groups[0]
	timestamp := s.Date
	duration := s.Duration

	ttdb.UpdateCmd(s.CacheDir, func(tx *bolt.Tx) error {
		if timestamp.IsZero() {
			return fmt.Errorf("you can't set the zero types")
		}

		gb, err := getOrCreateBucketConditionally(tx, group.String(), duration.IsZero())
		if gb == nil || err != nil {
			return err
		}

		rb, err := getOrCreateBucketConditionally(gb, "rec", duration.IsZero())
		if rb == nil || err != nil {
			return err
		}

		date_key := timestamp.String()
		ttdb.SetSeconds(rb, date_key, duration)

		return nil
	})
}

func AggFunc(s *State) {
	/*
		    group := s.Groups[0]
		    beg_date := s.BeginDate.String()
		    end_date := s.EndDate.String()
			var secs types.DaySeconds

			ttdb.ViewCmd(s.CacheDir, func(tx *bolt.Tx) error {
				m := getDateMap(tx, group, beg_date, end_date)
				for _, v := range m {
					secs += v
				}

				return nil
			})

			fmt.Printf("%s\n", secs.String())

	*/
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

func getListOfGroups(dir string) []string {
    groups := []string{}
	filepath.Walk(dir, func(path string, info os.FileInfo, e error) error {
        if info == nil {
			return filepath.SkipDir
        }

        group_name := getRelWithPanic(dir, path)
        if info.IsDir() {
            if group_name != "." && !types.IsValidGroupFolder(group_name) {
                return filepath.SkipDir
            }
		} else {
    		if types.IsValidGroupFile(group_name) {
                group_cleaned_name := types.CreateGroupFromString(group_name).String()
                groups = append(groups, group_cleaned_name)
    		}
        }

		return nil
	})

	return groups
}

func LsFunc(s *State) {
    dirs := []string{}
    if len(s.Groups) == 0 {
        dirs = append(dirs, "")
    } else {
        for _, group := range s.Groups {
            dirs = append(dirs, group.String())
        }
    }

    visited_groups := map[string]bool{}
    for _, dir := range dirs {
        groups := getListOfGroups(filepath.Join(s.DataDir, dir))
        for _, group := range groups {
            group_with_path := filepath.Join(dir, group)
            if _, exists := visited_groups[group_with_path]; !exists {
                visited_groups[group_with_path] = true
                fmt.Println(group_with_path)
            }
        }
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
		ttfile.AddTimeout(filepath.Join(s.DataDir, group.Filename()), *timestamp_to_write, *seconds_to_write)
	}
}
