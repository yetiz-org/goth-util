package xtype

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var str = String(strings.Repeat("string", 100))

func TestBytes(t *testing.T) {
	bytes := Bytes(str.Bytes())
	assert.Equal(t, str.String(), bytes.String())
}
