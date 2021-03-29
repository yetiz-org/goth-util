package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNumber(t *testing.T) {
	assert.True(t, IsDigits("123"))
	assert.True(t, IsDigits("0123"))
	assert.False(t, IsDigits("4.78"))
	assert.False(t, IsDigits("+123"))
}
