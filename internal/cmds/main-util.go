package cmds

import (
	"strings"
	"time"
)

var DATE_FORMAT_STRING string = "2006-01-02"

func clean_group(group string) string {
	fields := strings.FieldsFunc(group, func(c rune) bool { return c == '/' })

	newFields := []string{}
	for _, v := range fields {
		for i, r := range v {
			if r != '.' {
				newFields = append(newFields, v[i:])
				break
			}
		}
	}

	group = strings.Join(newFields, "/")
	return group
}

func formatTimestamp(timestamp time.Time) string {
	if timestamp.IsZero() {
		return ""
	}
	return timestamp.Format(DATE_FORMAT_STRING)
}
