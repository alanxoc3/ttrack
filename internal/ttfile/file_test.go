package ttfile_test

import (
	"os"
	"testing"

/*
	"github.com/alanxoc3/ttrack/internal/ttfile"
	"github.com/alanxoc3/ttrack/internal/types"
	"github.com/stretchr/testify/assert"
	*/

	"io/ioutil"
)

// Creates a temp directory then removes it after the test.
func execTest(t *testing.T, testFunc func(string)) {
	dir, err := ioutil.TempDir("", "ttrack-test-")
	defer os.RemoveAll(dir)
	if err != nil {
		t.Fatalf("%s", err)
	}
	testFunc(dir)
}

/*
func TestSeconds(t *testing.T) {
	date, _ := types.CreateDateFromString("2021-12-31")
	second := types.CreateSecondsFromString("1s")
	expected_second := types.CreateSecondsFromString("2s")
	execTest(t, func(dir string) {
		ttfile.AddTimeout(dir+"/file", *date, second)
		ttfile.AddTimeout(dir+"/file", *date, second)
		m := ttfile.GetDateSeconds(dir + "/file")
		assert.Equal(t, map[types.Date]types.DaySeconds{
			*date: expected_second,
		}, m)
	})
}
*/
