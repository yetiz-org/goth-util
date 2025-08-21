package value

import (
	"testing"
)

// BenchmarkJsonMarshal tests the performance of JsonMarshal function
func BenchmarkJsonMarshal(b *testing.B) {
	testCases := []interface{}{
		"simple string",
		123,
		123.456,
		true,
		map[string]interface{}{
			"name":  "test",
			"value": 42,
		},
		[]interface{}{1, 2, 3, 4, 5},
		struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}{"Alice", 30},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			_ = JsonMarshal(tc)
		}
	}
}

// BenchmarkJsonMarshalString tests JsonMarshal performance specifically for strings
func BenchmarkJsonMarshalString(b *testing.B) {
	str := "Hello, World! This is a test string for benchmarking."
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = JsonMarshal(str)
	}
}

// BenchmarkJsonMarshalInt tests JsonMarshal performance specifically for integers
func BenchmarkJsonMarshalInt(b *testing.B) {
	num := 42
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = JsonMarshal(num)
	}
}

// BenchmarkJsonMarshalComplex tests JsonMarshal performance for complex types
func BenchmarkJsonMarshalComplex(b *testing.B) {
	data := map[string]interface{}{
		"users": []map[string]interface{}{
			{"id": 1, "name": "Alice", "active": true},
			{"id": 2, "name": "Bob", "active": false},
		},
		"total": 2,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = JsonMarshal(data)
	}
}

// BenchmarkCopy tests the performance of Copy function
func BenchmarkCopy(b *testing.B) {
	type TestStruct struct {
		Name   string
		Age    int
		Active bool
		Score  float64
	}
	
	src := TestStruct{
		Name:   "Alice",
		Age:    30,
		Active: true,
		Score:  95.5,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var dst TestStruct
		Copy(src, &dst)
	}
}

// BenchmarkCopyLargeStruct tests Copy performance with larger structures
func BenchmarkCopyLargeStruct(b *testing.B) {
	type LargeStruct struct {
		Field1  string
		Field2  int
		Field3  bool
		Field4  float64
		Field5  string
		Field6  int64
		Field7  bool
		Field8  float32
		Field9  string
		Field10 int
	}
	
	src := LargeStruct{
		Field1:  "test1",
		Field2:  42,
		Field3:  true,
		Field4:  3.14,
		Field5:  "test5",
		Field6:  1234567890,
		Field7:  false,
		Field8:  2.71,
		Field9:  "test9",
		Field10: 100,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var dst LargeStruct
		Copy(src, &dst)
	}
}

// BenchmarkCast tests the performance of Cast function
func BenchmarkCast(b *testing.B) {
	value := "test value"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := Cast[string](value)
		_ = result
	}
}
