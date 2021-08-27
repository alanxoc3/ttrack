package types

import (
	"bufio"
	"strings"
)

const FILE_ENDING = ".tt"

type Group struct {
	val string
}

func IsValidGroupFolder(path string) bool {
    return len(path) > 0 && path == CreateGroupFromString(path).String()
}

func IsValidGroupFile(path string) bool {
    return len(path) > 0 && path == CreateGroupFromString(path).Filename()
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

// If oldest_ancestor isn't actually an ancestor, an empty list will be returned.
// Use the empty group to get all ancestors.
func (group Group) GetAncestors(oldest_ancestor Group) []Group {
    oldest_ancestor_dir := oldest_ancestor.val
    if len(oldest_ancestor_dir) > 0 { oldest_ancestor_dir = oldest_ancestor_dir + "/" }
    if oldest_ancestor_dir != "/" && !strings.HasPrefix(group.val, oldest_ancestor_dir) { return []Group{} }

    group_without_oldest_ancestor := strings.Replace(group.val, oldest_ancestor_dir, "", 1)
    dirs := strings.Split(group_without_oldest_ancestor, "/")
    ancestors := make([]Group, 0, len(dirs))
    for _, dir := range dirs {
        ancestors = append(ancestors, CreateGroupFromString(oldest_ancestor_dir + dir))
        oldest_ancestor_dir = oldest_ancestor_dir + dir + "/"
    }

    return ancestors
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
