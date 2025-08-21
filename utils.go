// Package kkutil provides common utility functions for various data types and operations.
package kkutil

import (
	"net"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// zeroTime represents the Unix epoch time (1970-01-01 00:00:00 UTC)
var zeroTime = time.Unix(0, 0)

// BytesFromUUID converts a UUID object to its binary representation as a byte slice.
// It uses the UUID's MarshalBinary method to get the raw bytes.
// Optimized version that directly accesses UUID bytes to avoid allocation.
func BytesFromUUID(uuidObj uuid.UUID) []byte {
	// UUID is already a [16]byte array, we can directly convert it
	// This avoids the allocation that MarshalBinary creates
	result := make([]byte, 16)
	copy(result, uuidObj[:])
	return result
}

// ZeroTimestamp returns the zero timestamp value (0) representing the Unix epoch.
// This is commonly used as a default or initial timestamp value.
func ZeroTimestamp() int64 {
	return 0
}

// MaxTimestamp returns the maximum timestamp value for 32-bit systems (2^32 - 1).
// This represents the maximum date that can be stored in a 32-bit Unix timestamp.
func MaxTimestamp() int64 {
	return 4294967295
}

// MaxTimestamp32u returns the maximum timestamp value as an unsigned integer.
// This is equivalent to MaxTimestamp but returns uint type instead of int64.
func MaxTimestamp32u() uint {
	return 4294967295
}

// UnixToTime converts a Unix timestamp (seconds since epoch) to a time.Time object.
// The nanosecond component is set to 0.
func UnixToTime(unixSecond int64) time.Time {
	return time.Unix(unixSecond, 0)
}

// IsInt checks if a string represents a valid integer.
// It attempts to parse the string as a 64-bit integer and returns true if successful.
func IsInt(str string) bool {
	if _, err := strconv.ParseInt(str, 10, 64); err != nil {
		return false
	}

	return true
}

// CastString attempts to cast an interface{} to a string pointer.
// Returns a pointer to the string if the cast is successful, nil otherwise.
func CastString(obj interface{}) *string {
	if str, ok := obj.(string); ok {
		return &str
	}

	return nil
}

// NowUInt returns the current Unix timestamp as an unsigned integer.
// This is a convenience function for getting the current time in Unix format.
func NowUInt() uint {
	return uint(time.Now().Unix())
}

// SplitRemoteAddr parses a network address string and extracts the IP and port components.
// It handles both IPv4 and IPv6 addresses with proper bracket notation for IPv6.
// Returns the parsed IP address and port string. If no port is found, port will be empty.
// Optimized version with reduced string operations and single-pass parsing.
func SplitRemoteAddr(addr string) (ip net.IP, port string) {
	if addr == "" {
		return nil, ""
	}

	// Fast path for IPv6 with brackets: [ip]:port
	if addr[0] == '[' {
		// Find the closing bracket
		closeBracket := -1
		for i := 1; i < len(addr); i++ {
			if addr[i] == ']' {
				closeBracket = i
				break
			}
		}
		
		if closeBracket == -1 {
			// Malformed - no closing bracket
			return nil, ""
		}
		
		// Extract IP part (without brackets)
		ipStr := addr[1:closeBracket]
		
		// Check if there's a port after the bracket
		if closeBracket+2 < len(addr) && addr[closeBracket+1] == ':' {
			port = addr[closeBracket+2:]
		}
		
		return net.ParseIP(ipStr), port
	}

	// Find the last colon for port separation
	lastColon := -1
	colonCount := 0
	for i := len(addr) - 1; i >= 0; i-- {
		if addr[i] == ':' {
			colonCount++
			if lastColon == -1 {
				lastColon = i
			}
		}
	}

	// No colon found - treat as IP only
	if lastColon == -1 {
		return net.ParseIP(addr), ""
	}

	// Multiple colons without brackets - likely IPv6 without port
	if colonCount > 1 {
		return net.ParseIP(addr), ""
	}

	// Single colon - IPv4 with port or hostname with port
	ipStr := addr[:lastColon]
	port = addr[lastColon+1:]
	return net.ParseIP(ipStr), port
}
