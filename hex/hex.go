// Package hex provides utilities for hexadecimal encoding and decoding operations.
// It wraps the standard library's encoding/hex package with additional error handling.
package hex

import "encoding/hex"

// EncodeToString encodes a byte slice to its hexadecimal string representation.
// This is a direct wrapper around the standard library's hex.EncodeToString function.
func EncodeToString(src []byte) string {
	return hex.EncodeToString(src)
}

// DecodeString decodes a hexadecimal string to its byte slice representation.
// Returns nil if the input string is not valid hexadecimal.
// This provides safer error handling compared to the standard library function.
// Optimized version with early validation to avoid unnecessary allocations.
func DecodeString(src string) []byte {
	// Fast path: check for empty string
	if src == "" {
		return []byte{}
	}
	
	// Fast path: check for odd length (invalid hex)
	if len(src)&1 == 1 {
		return nil
	}
	
	// Fast path: validate hex characters before calling standard library
	// This avoids allocation in the standard library for invalid inputs
	for i := 0; i < len(src); i++ {
		c := src[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return nil
		}
	}
	
	// Now we know it's valid, safe to call standard library
	if bytes, e := hex.DecodeString(src); e != nil {
		return nil
	} else {
		return bytes
	}
}
