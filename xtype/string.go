package xtype

import (
	"strings"
	"unsafe"
)

// String is an extended string type with additional utility methods.
// It provides convenient methods for common string operations.
type String string

// Bytes converts the string to a byte slice using unsafe pointer casting.
// This is more efficient than []byte() conversion as it avoids copying.
// Warning: The returned byte slice shares memory with the original string.
func (s String) Bytes() []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

// String returns the underlying string value.
// This method provides a consistent interface for accessing the raw string.
func (s String) String() string {
	return string(s)
}

// Upper converts the string to uppercase.
// This is a convenience method that wraps strings.ToUpper.
func (s String) Upper() string {
	return strings.ToUpper(string(s))
}

// Lower converts the string to lowercase.
// This is a convenience method that wraps strings.ToLower.
func (s String) Lower() string {
	return strings.ToLower(string(s))
}

// Short creates an abbreviated version of the string using the default separator "-".
// It takes the first character of each segment split by hyphens.
// For example, "hello-world-test" becomes "HWT".
func (s String) Short() string {
	return s.ShortBySign("-")
}

// ShortBySign creates an abbreviated version of the string using a custom separator.
// It splits the string by the provided sign, converts to uppercase, and takes
// the first character of each non-empty segment to form an abbreviation.
// For example, with sign=".", "hello.world.test" becomes "HWT".
// Optimized version using strings.Builder to reduce memory allocations.
func (s String) ShortBySign(sign string) string {
	upperStr := strings.ToUpper(string(s))
	segments := strings.Split(upperStr, sign)
	
	// Pre-allocate builder with estimated capacity
	var sb strings.Builder
	sb.Grow(len(segments)) // Each segment contributes at most 1 character
	
	for _, sp := range segments {
		if len(sp) > 0 {
			sb.WriteByte(sp[0])
		}
	}

	return sb.String()
}
