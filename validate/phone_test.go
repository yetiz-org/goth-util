package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerritoryCode(t *testing.T) {
	assert.True(t, IsValidPhoneTerritoryCode("tw"))
	assert.True(t, IsValidPhoneTerritoryCode("TW"))

	assert.False(t, IsValidPhoneTerritoryCode("TWN"))
	assert.False(t, IsValidPhoneTerritoryCode("AA"))
}

func TestCodeRegionEqual(t *testing.T) {
	assert.True(t, IsCodeRegionEqual("886", "TW"))
	assert.True(t, IsCodeRegionEqual("886", "tw"))
	assert.True(t, IsCodeRegionEqual("1", "US"))
	assert.True(t, IsCodeRegionEqual("1", "CA"))

	assert.False(t, IsCodeRegionEqual("", ""))
	assert.False(t, IsCodeRegionEqual("", "TW"))
	assert.False(t, IsCodeRegionEqual("886", ""))
	assert.False(t, IsCodeRegionEqual("TW", "TW"))
	assert.False(t, IsCodeRegionEqual("12345", "TW"))
	assert.False(t, IsCodeRegionEqual("86", "TW"))
}
