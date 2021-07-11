package types_test

import (
	"fmt"
	"testing"

	"github.com/alanxoc3/ttrack/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestCreateGroupFromString(t *testing.T) {
	var testVals = []struct {
		expectedValidFile bool
		expectedValidFolder bool
		expectedFile string
		expectedString string
		test     string
	}{
		{false, false, "", "", ""},
		{false, false, "", "", "/"},
		{false, false, "", "", "////////////"},
		{false, false, "", "", "/ /\n/\t/\n\t\r /"},
		{false, false, "", "", "/./../.../..../.....tt/"},
		{false, false, "", "", "/./../.../..../.....tt/"},
		{false, false, "", "", ".tt"},
		{false, false, "", "", "...tt"},
		{false, false, "a.tt", "a", " a "},
		{false, false, "a.tt", "a", " a.tt "},
		{false, false, "a.tt", "a", "/a/"},
		{false, false, "a.tt", "a", " / a / "},
		{true, false, "a..tt", "a.", "a..tt"},
		{false, false, "a/b.tt", "a/b", "/a/b/"},
		{false, false, "a/b.tt", "a/b", " / a / b / "},
		{false, true, "a/b.tt", "a/b", "a/b"},
		{false, false, "a/b.tt", "a/b", "a.tt/b.tt"},
		{false, false, "a/b.tt", "a/b", ".a/.b.tt"},
		{false, false, "a/b.tt", "a/b", ".....a/........b.tt"},
		{false, false, "a/b.tt", "a/b", "  .....a.tt.tt.tt.tt \t/ \n........b.tt.tt.tt.tt.tt"},
		{false, false, "a/b/c.tt", "a/b/c", ".a/.b/c.tt"},
	}
	for _, v := range testVals {
		t.Run(fmt.Sprintf("test-%s", v.test), func(t *testing.T) {
    		actual := types.CreateGroupFromString(v.test)
			assert.Equal(t, v.expectedValidFolder, types.IsValidGroupFolder(v.test))
			assert.Equal(t, v.expectedValidFile, types.IsValidGroupFile(v.test))
			assert.Equal(t, v.expectedString, actual.String())
			assert.Equal(t, v.expectedFile, actual.Filename())
		})
	}

}
