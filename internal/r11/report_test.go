package r11

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExcludeByOuterSymbolOnly(t *testing.T) {
	want := true
	got := excludeByOuterSymbolOnly("8107", "8107903528923")
	assert.Equal(t, want, got)
	want = false
	got = excludeByOuterSymbolOnly("8107", "81077903528923")
	assert.Equal(t, want, got)
}
