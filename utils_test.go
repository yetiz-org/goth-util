package kkutil

import (
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBytesFromUUID(t *testing.T) {
	// Test with a known UUID
	uuidStr := "550e8400-e29b-41d4-a716-446655440000"
	uuidObj, err := uuid.Parse(uuidStr)
	assert.NoError(t, err)

	bytes := BytesFromUUID(uuidObj)
	assert.Equal(t, 16, len(bytes))

	// Test that the bytes can be converted back to UUID
	newUUID, err := uuid.FromBytes(bytes)
	assert.NoError(t, err)
	assert.Equal(t, uuidObj, newUUID)

	// Test with random UUID
	randomUUID := uuid.New()
	randomBytes := BytesFromUUID(randomUUID)
	assert.Equal(t, 16, len(randomBytes))

	reconstructedUUID, err := uuid.FromBytes(randomBytes)
	assert.NoError(t, err)
	assert.Equal(t, randomUUID, reconstructedUUID)
}

func TestZeroTimestamp(t *testing.T) {
	assert.Equal(t, int64(0), ZeroTimestamp())
}

func TestMaxTimestamp(t *testing.T) {
	expected := int64(4294967295)
	assert.Equal(t, expected, MaxTimestamp())
}

func TestMaxTimestamp32u(t *testing.T) {
	expected := uint(4294967295)
	assert.Equal(t, expected, MaxTimestamp32u())
}

func TestUnixToTime(t *testing.T) {
	// Test with zero timestamp
	zeroTime := UnixToTime(0)
	expected := time.Unix(0, 0)
	assert.Equal(t, expected, zeroTime)

	// Test with specific timestamp
	timestamp := int64(1640995200) // 2022-01-01 00:00:00 UTC
	result := UnixToTime(timestamp)
	expected = time.Unix(timestamp, 0)
	assert.Equal(t, expected, result)

	// Test with negative timestamp
	negativeTimestamp := int64(-86400) // 1 day before epoch
	result = UnixToTime(negativeTimestamp)
	expected = time.Unix(negativeTimestamp, 0)
	assert.Equal(t, expected, result)
}

func TestIsInt(t *testing.T) {
	// Test valid integers
	assert.True(t, IsInt("123"))
	assert.True(t, IsInt("-456"))
	assert.True(t, IsInt("0"))
	assert.True(t, IsInt("9223372036854775807")) // max int64

	// Test invalid strings
	assert.False(t, IsInt("123.45"))
	assert.False(t, IsInt("abc"))
	assert.False(t, IsInt("12a34"))
	assert.False(t, IsInt(""))
	assert.False(t, IsInt(" 123"))
	assert.False(t, IsInt("123 "))
	assert.False(t, IsInt("1.0"))
}

func TestCastString(t *testing.T) {
	// Test successful cast
	str := "hello"
	result := CastString(str)
	assert.NotNil(t, result)
	assert.Equal(t, str, *result)

	// Test failed cast - integer
	result = CastString(123)
	assert.Nil(t, result)

	// Test failed cast - nil
	result = CastString(nil)
	assert.Nil(t, result)

	// Test failed cast - other types
	result = CastString([]string{"hello"})
	assert.Nil(t, result)

	result = CastString(map[string]int{"key": 1})
	assert.Nil(t, result)
}

func TestNowUInt(t *testing.T) {
	before := uint(time.Now().Unix())
	result := NowUInt()
	after := uint(time.Now().Unix())

	// The result should be between before and after (allowing for time passage)
	assert.True(t, result >= before)
	assert.True(t, result <= after)
}

func TestSplitRemoteAddr(t *testing.T) {
	// Test IPv4 with port
	ip, port := SplitRemoteAddr("192.168.1.1:8080")
	assert.Equal(t, net.ParseIP("192.168.1.1"), ip)
	assert.Equal(t, "8080", port)

	// Test IPv4 without port
	ip, port = SplitRemoteAddr("192.168.1.1")
	assert.Equal(t, net.ParseIP("192.168.1.1"), ip)
	assert.Equal(t, "", port)

	// Test IPv6 with port and brackets
	ip, port = SplitRemoteAddr("[2001:db8::1]:8080")
	assert.Equal(t, net.ParseIP("2001:db8::1"), ip)
	assert.Equal(t, "8080", port)

	// Test IPv6 without port
	ip, port = SplitRemoteAddr("2001:db8::1")
	assert.Equal(t, net.ParseIP("2001:db8::1"), ip)
	assert.Equal(t, "", port)

	// Test IPv6 localhost with port
	ip, port = SplitRemoteAddr("[::1]:9000")
	assert.Equal(t, net.ParseIP("::1"), ip)
	assert.Equal(t, "9000", port)

	// Test localhost without port
	ip, port = SplitRemoteAddr("127.0.0.1")
	assert.Equal(t, net.ParseIP("127.0.0.1"), ip)
	assert.Equal(t, "", port)

	// Test edge cases
	ip, port = SplitRemoteAddr("localhost:3000")
	assert.Nil(t, ip) // localhost is not a valid IP
	assert.Equal(t, "3000", port) // but port should still be extracted

	// Test empty string
	ip, port = SplitRemoteAddr("")
	assert.Nil(t, ip)
	assert.Equal(t, "", port)

	// Test malformed IPv6
	ip, port = SplitRemoteAddr("[2001:db8::1:8080")
	assert.Nil(t, ip) // malformed IPv6 should return nil
	assert.Equal(t, "", port)

	// Test IPv6 with multiple colons without brackets
	ip, port = SplitRemoteAddr("fe80::1%lo0")
	assert.Equal(t, net.ParseIP("fe80::1%lo0"), ip)
	assert.Equal(t, "", port)
}
