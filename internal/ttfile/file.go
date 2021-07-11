package ttfile

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alanxoc3/ttrack/internal/types"
)

func insertOrAdd(dateToSeconds map[types.Date]types.DaySeconds, day types.Date, newDurationFunc func(types.DaySeconds)types.DaySeconds) {
	if day.IsZero() {
		return
	}

	if val, ok := dateToSeconds[day]; ok {
		dateToSeconds[day] = newDurationFunc(val)
	} else {
		dateToSeconds[day] = newDurationFunc(types.CreateSecondsFromUint32(0))
	}
}

func GetDateSeconds(filename string) map[types.Date]types.DaySeconds {
	dateToSeconds := map[types.Date]types.DaySeconds{}

	if f, err := os.Open(filename); err == nil {
		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			line := scanner.Text()
			tokens := strings.Split(line, ":")
			if len(tokens) == 2 {
				// TODO: Log errors?
				d, _ := types.CreateDateFromString(strings.TrimSpace(tokens[0]))
				s := types.CreateSecondsFromString(strings.TrimSpace(tokens[1]))
				if d != nil {
					insertOrAdd(dateToSeconds, *d, func(ds types.DaySeconds)types.DaySeconds {
                        return ds.Add(s)
					})
				}
			} else {
				// TODO: Log errors?
			}
		}

		f.Close()

		// If the file doesn't exist, continue with writing.
		// But if the error is something else, we should exit the program.
	} else if !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}

	return dateToSeconds
}

func ModifyTime(filename string, insertion_date types.Date, newDurationFunc func(types.DaySeconds)types.DaySeconds) {
	date_list := types.DateList{}
	lines := GetDateSeconds(filename)

	insertOrAdd(lines, insertion_date, newDurationFunc)

	for k := range lines {
		date_list = append(date_list, k)
	}

	sort.Sort(date_list)

	output_str := ""
	for _, k := range date_list {
		v := lines[k]
		output_str += k.String() + ": " + v.String() + "\n"
	}

    // Create dir with execute for cd, if it doesn't already exist.
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

    // Create file with rw for user and r for everyone else (before umask).
	if err := os.WriteFile(filename, []byte(output_str), 0644); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}
}
