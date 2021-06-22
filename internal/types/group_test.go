package types_test

import (
	"fmt"
	"testing"

	"github.com/alanxoc3/ttrack/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestCreateGroupFromString(t *testing.T) {
	var testVals = []struct {
		expectedFile string
		expected string
		test     string
	}{
		{"", "", ""},
		{"", "", "/"},
		{"", "", "////////////"},
		{"", "", "/ /\n/\t/\n\t\r /"},
		{"", "", "/./../.../..../.....tt/"},
		{"", "", "/./../.../..../.....tt/"},
		{"", "", ".tt"},
		{"", "", "...tt"},
		{"a.tt", "a", " a "},
		{"a.tt", "a", " a.tt "},
		{"a.tt", "a", "/a/"},
		{"a.tt", "a", " / a / "},
		{"a..tt", "a.", "a..tt"},
		{"a/b.tt", "a/b", "/a/b/"},
		{"a/b.tt", "a/b", " / a / b / "},
		{"a/b.tt", "a/b", "a/b"},
		{"a/b.tt", "a/b", "a.tt/b.tt"},
		{"a/b.tt", "a/b", ".a/.b.tt"},
		{"a/b.tt", "a/b", ".....a/........b.tt"},
		{"a/b.tt", "a/b", "  .....a.tt.tt.tt.tt \t/ \n........b.tt.tt.tt.tt.tt"},
		{"a/b/c.tt", "a/b/c", ".a/.b/c.tt"},
	}
	for _, v := range testVals {
		actual := types.CreateGroupFromString(v.test)
		t.Run(fmt.Sprintf("test-%s", v.expected), func(t *testing.T) {
			assert.Equal(t, v.expected, actual.String())
			assert.Equal(t, v.expectedFile, actual.Filename())
		})
	}

}
