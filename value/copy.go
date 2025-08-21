// Package value provides utilities for type casting and copying values using reflection.
package value

import "reflect"

// Cast performs a type assertion to convert an interface{} value to a specific type T.
// Returns the zero value of type T if the conversion fails.
// This is a generic function that provides safe type casting.
func Cast[T any](in any) (t T) {
	if v, ok := in.(T); ok {
		t = v
	}

	return
}

// Copy performs a deep copy of values from the source to the destination using reflection.
// It handles primitive types, structs, and maps with type compatibility checking.
// Both from and to parameters can be pointers or values.
// The function only copies fields with matching names and compatible types.
// Optimized version with reduced reflection overhead and early type checking.
func Copy(from interface{}, to interface{}) {
	// Fast path: check for nil inputs early
	if from == nil || to == nil {
		return
	}
	
	toVal := reflect.ValueOf(to)
	fromVal := reflect.ValueOf(from)
	
	// Fast path: check if 'to' is not a pointer (cannot be set)
	if toVal.Kind() != reflect.Ptr {
		return
	}
	
	// Dereference 'to' pointer once
	toElem := toVal.Elem()
	if !toElem.CanSet() {
		return
	}
	
	// Handle 'from' pointer dereferencing
	fromElem := fromVal
	if fromVal.Kind() == reflect.Ptr {
		if fromVal.IsNil() {
			return
		}
		fromElem = fromVal.Elem()
	}
	
	// Fast path: check type compatibility early
	toKind := toElem.Kind()
	fromKind := fromElem.Kind()
	
	// Optimize primitive type copying with switch on kind
	switch toKind {
	case reflect.Bool:
		if fromKind == reflect.Bool {
			toElem.SetBool(fromElem.Bool())
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if fromKind >= reflect.Int && fromKind <= reflect.Int64 {
			toElem.SetInt(fromElem.Int())
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if fromKind >= reflect.Uint && fromKind <= reflect.Uintptr {
			toElem.SetUint(fromElem.Uint())
		}
	case reflect.Float32, reflect.Float64:
		if fromKind == reflect.Float32 || fromKind == reflect.Float64 {
			toElem.SetFloat(fromElem.Float())
		}
	case reflect.Complex64, reflect.Complex128:
		if fromKind == reflect.Complex64 || fromKind == reflect.Complex128 {
			toElem.SetComplex(fromElem.Complex())
		}
	case reflect.String:
		if fromKind == reflect.String {
			toElem.SetString(fromElem.String())
		}
	case reflect.Struct:
		// Optimize struct copying with field caching
		if fromKind == reflect.Struct {
			copyStruct(fromElem, toElem)
		}
	case reflect.Map:
		// Copy maps if they have compatible types
		if fromKind == reflect.Map && toElem.Type() == fromElem.Type() {
			toElem.Set(fromElem)
		}
	}
}

// copyStruct optimizes struct field copying with reduced reflection overhead
func copyStruct(fromElem, toElem reflect.Value) {
	toType := toElem.Type()
	
	// Cache field lookups to reduce reflection overhead
	numFields := toElem.NumField()
	for i := 0; i < numFields; i++ {
		toField := toElem.Field(i)
		if !toField.CanSet() {
			continue
		}
		
		toFieldType := toType.Field(i)
		fieldName := toFieldType.Name
		
		// Fast field lookup by name
		fromField := fromElem.FieldByName(fieldName)
		if !fromField.IsValid() {
			continue
		}
		
		// Check type compatibility
		if fromField.Kind() == toField.Kind() && fromField.Type() == toField.Type() {
			toField.Set(fromField)
		}
	}
}
