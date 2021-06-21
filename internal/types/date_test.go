package types_test

import (
	"testing"
	"time"

	"github.com/alanxoc3/ttrack/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestCreateFromString(t *testing.T) {
	d, err := types.CreateDateFromString("2021-01-01")
	assert.Nil(t, err)
	assert.Equal(t, "2021-01-01", d.String())
}

func TestCreateFromStringError(t *testing.T) {
	d, err := types.CreateDateFromString("2021-01-1")
	assert.NotNil(t, err)
	assert.Nil(t, d)
}

func TestCreateFromTimeWithZeroTime(t *testing.T) {
	d := types.CreateDateFromTime(time.Time{})
	assert.NotNil(t, d)
	assert.Equal(t, "0001-01-01", d.String())
	assert.True(t, d.IsZero())
}
