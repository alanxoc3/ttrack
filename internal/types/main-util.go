package types

import (
	"strings"
)

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
