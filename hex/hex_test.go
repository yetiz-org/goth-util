package hex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeToString(t *testing.T) {
	// Test basic encoding
	input := []byte("hello")
	expected := "68656c6c6f"
	result := EncodeToString(input)
	assert.Equal(t, expected, result)

	// Test empty byte slice
	result = EncodeToString([]byte{})
	assert.Equal(t, "", result)

	// Test single byte
	result = EncodeToString([]byte{255})
	assert.Equal(t, "ff", result)

	// Test multiple bytes with various values
	input = []byte{0, 15, 16, 255}
	expected = "000f10ff"
	result = EncodeToString(input)
	assert.Equal(t, expected, result)

	// Test nil input
	result = EncodeToString(nil)
	assert.Equal(t, "", result)
}

func TestDecodeString(t *testing.T) {
	// Test basic decoding
	input := "68656c6c6f"
	expected := []byte("hello")
	result := DecodeString(input)
	assert.Equal(t, expected, result)

	// Test empty string
	result = DecodeString("")
	assert.Equal(t, []byte{}, result)

	// Test single byte
	result = DecodeString("ff")
	assert.Equal(t, []byte{255}, result)

	// Test multiple bytes
	input = "000f10ff"
	expected = []byte{0, 15, 16, 255}
	result = DecodeString(input)
	assert.Equal(t, expected, result)

	// Test uppercase hex
	input = "68656C6C6F"
	expected = []byte("hello")
	result = DecodeString(input)
	assert.Equal(t, expected, result)

	// Test invalid hex strings
	result = DecodeString("xyz")
	assert.Nil(t, result)

	result = DecodeString("1")
	assert.Nil(t, result) // odd length

	result = DecodeString("1g")
	assert.Nil(t, result) // invalid character

	result = DecodeString("hello")
	assert.Nil(t, result) // non-hex string
}

func TestRoundTrip(t *testing.T) {
	// Test encode -> decode round trip
	original := []byte("Hello, World! 123 !@#$%^&*()")
	encoded := EncodeToString(original)
	decoded := DecodeString(encoded)
	assert.Equal(t, original, decoded)

	// Test with binary data
	original = []byte{0, 1, 2, 3, 255, 254, 253}
	encoded = EncodeToString(original)
	decoded = DecodeString(encoded)
	assert.Equal(t, original, decoded)

	// Test with random bytes
	original = make([]byte, 256)
	for i := 0; i < 256; i++ {
		original[i] = byte(i)
	}
	encoded = EncodeToString(original)
	decoded = DecodeString(encoded)
	assert.Equal(t, original, decoded)
}

func TestDecodeStringErrorHandling(t *testing.T) {
	// Test various error conditions that should return nil
	testCases := []string{
		"zz",           // invalid hex characters
		"1",            // odd length
		"1z",           // mixed valid/invalid
		"hello world",  // spaces and invalid chars
		"12 34",        // spaces in middle
		"G0",           // invalid uppercase
		"0G",           // invalid at end
	}

	for _, testCase := range testCases {
		result := DecodeString(testCase)
		assert.Nil(t, result, "Expected nil for input: %s", testCase)
	}
}
