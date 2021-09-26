package cmds

import (
	"errors"
	"io/fs"
	"os"
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
	if beg_ts.IsZero() || end_ts.IsZero() || timeout.IsZero() || now.Before(end_ts) || end_ts.Before(beg_ts) {
		return types.DaySeconds{}, false
	} else if now.Sub(end_ts).Seconds() > float64(timeout.GetAsUint32()) {
		return types.CreateSecondsFromUint32(uint32(end_ts.Sub(beg_ts).Seconds())).Add(timeout), true
	} else {
		return types.CreateSecondsFromUint32(uint32(now.Sub(beg_ts).Seconds())), false
	}
}

func expandGroup(b *bolt.Bucket) (time.Time, time.Time, types.DaySeconds) {
	return ttdb.GetTimestamp(b, "beg"), ttdb.GetTimestamp(b, "end"), ttdb.GetSeconds(b, "out")
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
		if !info.IsDir() && types.IsValidGroupFile(group_name) {
			groups = append(groups, types.CreateGroupFromString(group_name).String())
		}

		if (strat == walk_recursive || strat == walk_level) && info.IsDir() && types.IsValidGroupFolder(group_name) {
			groups = append(groups, types.CreateGroupFromString(group_name).String())
		}

		if info.IsDir() && group_name != "." && (strat == walk_level || !types.IsValidGroupFolder(group_name)) {
			return filepath.SkipDir
		} else {
			return nil
		}
	})

	return groups
}

func walkThroughGroups(cache_dir, data_dir string, groupdirs []types.Group, strat walkstrategy, cached, stored bool) map[types.Group]bool {
	visited_groups := map[types.Group]bool{}

	if cached && (strat == walk_level || strat == walk_recursive) {
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
	}

	if stored {
		for _, groupdir := range groupdirs {
			// Check for the group itself. It could be a folder or a .tt file.
			if !groupdir.IsZero() {
				_, file_err := os.Stat(filepath.Join(data_dir, groupdir.String()))
				if !errors.Is(file_err, fs.ErrNotExist) {
					if _, exists := visited_groups[groupdir]; !exists {
						visited_groups[groupdir] = true
					}
				}

				_, folder_err := os.Stat(filepath.Join(data_dir, groupdir.Filename()))
				if (strat == walk_level || strat == walk_recursive) && !errors.Is(folder_err, fs.ErrNotExist) {
					if _, exists := visited_groups[groupdir]; !exists {
						visited_groups[groupdir] = true
					}
				}
			}

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
	}

	return visited_groups
}
