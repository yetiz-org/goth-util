package hex

import (
	"crypto/rand"
	"testing"
)

// Benchmark for EncodeToString function
func BenchmarkEncodeToString(b *testing.B) {
	// Create test data of different sizes
	testSizes := []int{16, 64, 256, 1024, 4096}
	testData := make([][]byte, len(testSizes))
	
	for i, size := range testSizes {
		testData[i] = make([]byte, size)
		rand.Read(testData[i])
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data := testData[i%len(testData)]
		EncodeToString(data)
	}
}

// Benchmark for DecodeString function
func BenchmarkDecodeString(b *testing.B) {
	// Create test hex strings of different sizes
	testSizes := []int{16, 64, 256, 1024, 4096}
	testStrings := make([]string, len(testSizes))
	
	for i, size := range testSizes {
		data := make([]byte, size)
		rand.Read(data)
		testStrings[i] = EncodeToString(data)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hexStr := testStrings[i%len(testStrings)]
		DecodeString(hexStr)
	}
}

// Benchmark for round-trip operations
func BenchmarkHexRoundTrip(b *testing.B) {
	testData := make([]byte, 256)
	rand.Read(testData)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoded := EncodeToString(testData)
		DecodeString(encoded)
	}
}
