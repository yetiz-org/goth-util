package validate

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPublicIP(t *testing.T) {
	assert.False(t, IsPublicIP(net.ParseIP("10.24.2.62")))
	assert.True(t, IsPublicIP(net.ParseIP("110.24.2.62")))
}
