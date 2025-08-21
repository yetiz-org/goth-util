// Package xtype provides extended type definitions with additional utility methods.
package xtype

import (
	"unsafe"

	"github.com/yetiz-org/goth-util/hex"
)

// Bytes is a byte slice type with additional utility methods.
// It provides convenient methods for common byte slice operations.
type Bytes []byte

// String converts the byte slice to a string using unsafe pointer casting.
// This is more efficient than string() conversion as it avoids copying.
// Warning: The returned string shares memory with the original byte slice.
func (s Bytes) String() string {
	return *(*string)(unsafe.Pointer(&s))
}

// Bytes returns the underlying byte slice.
// This method provides a consistent interface for accessing the raw bytes.
func (s Bytes) Bytes() []byte {
	return s
}

// Hex returns the hexadecimal string representation of the byte slice.
// This is a convenience method that wraps the hex package's EncodeToString function.
func (s Bytes) Hex() string {
	return hex.EncodeToString(s)
}
