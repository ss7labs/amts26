package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundingToMinute(t *testing.T) {
	want := 60
	durations := []int{4, 6, 10, 11, 59, 1, 0}
	var count int
	for _, sec := range durations {
		res, err := CalculateBillingDuration(sec, 60, 0, 0)
		if err != nil {
			continue
		}
		assert.Equal(t, want, res)
		count++
	}
	assert.Equal(t, len(durations)-1, count)
}
