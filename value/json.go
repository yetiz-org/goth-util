package value

import (
	"encoding/json"
	"strconv"
)

// JsonMarshal converts an object to its JSON string representation.
// Returns an empty string if the marshaling fails.
// This provides a safe way to convert objects to JSON without handling errors explicitly.
// Optimized version with early type checking and reduced allocations.
func JsonMarshal(obj interface{}) string {
	// Fast path for nil
	if obj == nil {
		return "null"
	}
	
	// Fast path for common simple types to avoid reflection overhead
	switch v := obj.(type) {
	case string:
		// Pre-allocate buffer with estimated size to avoid reallocations
		buf := make([]byte, 0, len(v)+2) // +2 for quotes
		buf = append(buf, '"')
		// Simple escape for basic cases, fall back to json.Marshal for complex ones
		needsEscape := false
		for i := 0; i < len(v); i++ {
			c := v[i]
			if c < 32 || c == '"' || c == '\\' || c > 126 {
				needsEscape = true
				break
			}
		}
		if !needsEscape {
			buf = append(buf, v...)
			buf = append(buf, '"')
			return string(buf)
		}
		// Fall through to standard marshaling for strings needing escaping
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		if v {
			return "true"
		}
		return "false"
	}
	
	// For complex types, use standard library
	if marshal, err := json.Marshal(obj); err != nil {
		return ""
	} else {
		return string(marshal)
	}
}
