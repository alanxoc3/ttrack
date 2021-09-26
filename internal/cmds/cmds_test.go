package cmds_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"io/ioutil"

	"github.com/alanxoc3/ttrack/internal/cmds"
	"github.com/alanxoc3/ttrack/internal/types"
	"github.com/stretchr/testify/assert"
)

// Creates a temp directory then removes it after the test.
func execTest(t *testing.T, testFunc func(string, string)) {
	dir, err := ioutil.TempDir("", "ttrack-test-")
	defer os.RemoveAll(dir)
	if err != nil {
		t.Fatalf("%s", err)
	}
	testFunc(dir+"/cache", dir+"/data")
}

func mkstate(cache, data, group string) *cmds.State {
    return &cmds.State{
        CacheDir: cache,
        DataDir: data,
        Groups: []types.Group{types.CreateGroupFromString(group)},
        Cached: true,
        Stored: true,
    }
}

func mkstate_durdate(cache, data, group, timestamp, duration string) *cmds.State {
    s := mkstate(cache, data, group)
	now, _ := time.Parse(types.DATE_FORMAT_STRING, timestamp)
    date, _ := types.CreateDateFromString(timestamp)
    s.Now = now
    s.Date = *date
    s.Duration = types.CreateSecondsFromString(duration)
    return s
}

func fileToStr(t *testing.T, path ...string) string {
    b, err := os.ReadFile(strings.Join(path, "/"))
    assert.Nil(t, err)
    return string(b)
}

func TestBasicRec(t *testing.T) {
	execTest(t, func(cache, data string) {
 	    cmds.RecFunc(mkstate_durdate(cache, data, "testing", "2021-01-01", "10s"))
 	    cmds.RecFunc(mkstate_durdate(cache, data, "testing", "2021-01-02", "0s"))
 	    str := cmds.AggFunc(mkstate(cache, data, "testing"))
 	    assert.Equal(t, "10s\n", str)
	})
}

func TestReset(t *testing.T) {
	execTest(t, func(cache, data string) {
 	    cmds.RecFunc(mkstate_durdate(cache, data, "testing", "2021-01-01", "10s"))
 	    preReset := cmds.LsFunc(&cmds.State{CacheDir: cache, DataDir: data, Cached: true})
 	    assert.Equal(t, "testing\n", preReset)

 	    cmds.ResetFunc(mkstate(cache, data, "testing"))
 	    postReset := cmds.LsFunc(&cmds.State{CacheDir: cache, DataDir: data, Cached: true})
 	    assert.Zero(t, postReset)
	})
}

func TestAggBeforeRecDate(t *testing.T) {
	execTest(t, func(cache, data string) {
 	    cmds.RecFunc(mkstate_durdate(cache, data, "testing", "2021-01-02", "1s"))
        aggstr := cmds.AggFunc(mkstate(cache, data, "testing"))
 	    assert.Equal(t, "0s\n", aggstr)
	})
}

func TestTidy(t *testing.T) {
	execTest(t, func(cache, data string) {
 	    cmds.RecFunc(mkstate_durdate(cache, data, "testing", "2021-01-01", "10s"))
 	    {
     	    b, err := os.ReadFile(data+"/testing.tt")
     	    assert.Zero(t, b)
     	    assert.NotNil(t, err)
 	    }

 	    cmds.TidyFunc(mkstate_durdate(cache, data, "testing", "2021-01-02", ""))
 	    assert.Equal(t, "2021-01-01: 10s\n", fileToStr(t, data, "testing.tt"))
	})
}

func TestTidyExistingGroup(t *testing.T) {
	execTest(t, func(cache, data string) {
 	    cmds.AddFunc(mkstate_durdate(cache, data, "testing", "2021-01-01", "10s"))
 	    cmds.RecFunc(mkstate_durdate(cache, data, "testing", "2021-01-01", "3s"))

 	    cmds.TidyFunc(mkstate_durdate(cache, data, "testing", "2021-01-02", ""))
 	    assert.Equal(t, "2021-01-01: 13s\n", fileToStr(t, data, "testing.tt"))
	})
}

func TestSub(t *testing.T) {
	execTest(t, func(cache, data string) {
 	    cmds.AddFunc(mkstate_durdate(cache, data, "testing", "2021-01-01", "10s"))
 	    cmds.SubFunc(mkstate_durdate(cache, data, "testing", "2021-01-01", "6s"))

 	    str := cmds.AggFunc(mkstate(cache, data, "testing"))
 	    assert.Equal(t, "4s\n", str)
	})
}

func TestLsRecursiveCache(t *testing.T) {
	execTest(t, func(cache, data string) {
 	    cmds.RecFunc(mkstate_durdate(cache, data, "a/b/c/d", "2021-01-01", "3s"))
 	    str := cmds.LsFunc(&cmds.State{CacheDir: cache, DataDir: data, Recursive: true, Cached: true})
 	    assert.Equal(t, "a\na/b\na/b/c\na/b/c/d\n", str)
	})
}

func TestLsRecursiveFile(t *testing.T) {
	execTest(t, func(cache, data string) {
 	    cmds.AddFunc(mkstate_durdate(cache, data, "a/b/c/d", "2021-01-01", "3s"))
 	    str := cmds.LsFunc(&cmds.State{CacheDir: cache, DataDir: data, Recursive: true, Stored: true})
 	    assert.Equal(t, "a\na/b\na/b/c\na/b/c/d\n", str)
	})
}
