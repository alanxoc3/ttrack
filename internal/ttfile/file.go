package ttfile

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/alanxoc3/ttrack/internal/date"
	"github.com/alanxoc3/ttrack/internal/seconds"
)

type dateList []date.Date
func (dl dateList) Len() int { return len(dl) }
func (dl dateList) Less(i, j int) bool { return dl[i].IsLessThan(dl[j]) }
func (dl dateList) Swap(i, j int) { dl[i], dl[j] = dl[j], dl[i] }

func AddTimeout(filename string, insertion_date date.Date, timeout seconds.Seconds) {
	lines := map[date.Date]seconds.Seconds{}
	date_list := dateList{}

	if f, err := os.Open(filename); err == nil {
		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			line := scanner.Text()
			tokens := strings.Split(line, ":")
			if len(tokens) == 2 {
				// TODO: Log errors?
				d, _ := date.CreateFromString(strings.TrimSpace(tokens[0]))
				s := seconds.CreateFromString(strings.TrimSpace(tokens[1]))
				if d != nil && !d.IsZero() && s != 0 {
					if val, ok := lines[*d]; ok {
						lines[*d] = val + s
					} else {
						lines[*d] = s
						date_list = append(date_list, *d)
					}
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

	if !insertion_date.IsZero() && timeout != 0 {
		if val, ok := lines[insertion_date]; ok {
			lines[insertion_date] = val + timeout
		} else {
			lines[insertion_date] = timeout
			date_list = append(date_list, insertion_date)
		}
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
