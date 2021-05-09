package seconds_test

import (
	"testing"
	"time"

	"github.com/alanxoc3/ttrack/internal/seconds"
	"github.com/stretchr/testify/assert"
)

func TestCreateFromDuration(t *testing.T) {
	dur, _ := time.ParseDuration("1m3s")
	sec := seconds.CreateFromDuration(dur)
	assert.Equal(t, seconds.Seconds(63), sec)
}
