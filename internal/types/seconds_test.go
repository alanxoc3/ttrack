package types_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/alanxoc3/ttrack/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestCreateFromDuration(t *testing.T) {
	dur, _ := time.ParseDuration("1m3s")
	sec := types.CreateSecondsFromDuration(dur)
	assert.Equal(t, uint32(63), sec.GetAsUint32())
}

func TestString(t *testing.T) {
	var testVals = []struct {
		expected string
		test     time.Duration
	}{
		{"0s", 0},
		{"1s", 1},
		{"59s", 59},
		{"1m0s", 60},
		{"1m1s", 61},
		{"1h0m0s", 60 * 60},
		{"1h0m1s", 60*60 + 1},
		{"23h0m0s", 60 * 60 * 23},
		{"24h0m0s", 60 * 60 * 24},
		{"24h0m0s", 60 * 60 * 25},
		{"24h0m0s", (1 << 31 - 1) << 1 + 1}, // Max uint32
	}
	for _, v := range testVals {
		seconds := types.CreateSecondsFromDuration(v.test * time.Second)
		t.Run(fmt.Sprintf("test-%s", v.expected), func(t *testing.T) {
			assert.Equal(t, v.expected, seconds.String())
		})
	}
}
