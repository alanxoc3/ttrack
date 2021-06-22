package types

import (
	"bufio"
	"strings"
)

const FILE_ENDING = ".tt"

type Group struct {
	val string
}

func CreateGroupFromString(group_str string) Group {
    // STEP 1: Divide string into directories.
	fields := strings.FieldsFunc(group_str, func(c rune) bool { return c == '/' })

	newFields := []string{}
	for _, field := range fields {
		field_str := ""

		// STEP 2: For each path, remove all whitespace.
		scanner := bufio.NewScanner(strings.NewReader(field))
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() { field_str += scanner.Text() }
        if err := scanner.Err(); err != nil { panic(err) }

		// STEP 3: For each path, remove the filetype suffix.
		for strings.HasSuffix(field_str, FILE_ENDING) { field_str = strings.TrimSuffix(field_str, FILE_ENDING) }

		// STEP 4: For each path, remove leading dots.
		for strings.HasPrefix(field_str, ".") { field_str = strings.TrimPrefix(field_str, ".") }

        // STEP 5: Add the path to the path list if not empty.
        if len(field_str) > 0 { newFields = append(newFields, field_str) }
	}

	return Group{strings.Join(newFields, "/")}
}

func (group Group) String() string {
	return group.val
}

func (group Group) IsZero() bool {
    return group.val == ""
}

func (group Group) Filename() string {
    if !group.IsZero() {
    	return group.val + FILE_ENDING
    }
    return ""
}
