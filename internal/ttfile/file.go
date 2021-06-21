package ttfile

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/alanxoc3/ttrack/internal/types"
)

type dateList []types.Date
func (dl dateList) Len() int { return len(dl) }
func (dl dateList) Less(i, j int) bool { return dl[i].IsLessThan(dl[j]) }
func (dl dateList) Swap(i, j int) { dl[i], dl[j] = dl[j], dl[i] }

type dateSeconds struct {
    Date types.Date
    Seconds types.Seconds
}

func insertOrAdd(dateToSeconds map[types.Date]types.Seconds, day types.Date, secs types.Seconds) {
	if day.IsZero() || secs == 0 { return }

	if val, ok := dateToSeconds[day]; ok {
		dateToSeconds[day] = (val + secs).CapAtOneDay()
	} else {
		dateToSeconds[day] = secs.CapAtOneDay()
	}
}

func GetDateSeconds(filename string) map[types.Date]types.Seconds {
	dateToSeconds := map[types.Date]types.Seconds{}

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
				if d != nil { insertOrAdd(dateToSeconds, *d, s) }
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

func AddTimeout(filename string, insertion_date types.Date, timeout types.Seconds) {
	date_list := dateList{}
    lines := GetDateSeconds(filename)

    insertOrAdd(lines, insertion_date, timeout)

    for k := range lines {
        date_list = append(date_list, k)
    }

	sort.Sort(date_list)

	output_str := ""
	for _, k := range date_list {
    	v := lines[k]
		output_str += k.String() + ": " + v.String() + "\n"
	}

	// Write to a file
	err := os.WriteFile(filename, []byte(output_str), 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}