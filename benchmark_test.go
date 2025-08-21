package kkutil

import (
	"testing"

	"github.com/google/uuid"
)

// Benchmark for SplitRemoteAddr function
func BenchmarkSplitRemoteAddr(b *testing.B) {
	testCases := []string{
		"192.168.1.1:8080",
		"[2001:db8::1]:9000",
		"127.0.0.1:3000",
		"localhost:8080",
		"fe80::1%lo0",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		addr := testCases[i%len(testCases)]
		SplitRemoteAddr(addr)
	}
}

// Benchmark for BytesFromUUID function
func BenchmarkBytesFromUUID(b *testing.B) {
	uuids := make([]uuid.UUID, 100)
	for i := range uuids {
		uuids[i] = uuid.New()
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BytesFromUUID(uuids[i%len(uuids)])
	}
}

// Benchmark for IsInt function
func BenchmarkIsInt(b *testing.B) {
	testStrings := []string{
		"123456",
		"not_a_number",
		"9223372036854775807",
		"-456789",
		"123.45",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsInt(testStrings[i%len(testStrings)])
	}
}

// Benchmark for CastString function
func BenchmarkCastString(b *testing.B) {
	testValues := []interface{}{
		"hello world",
		42,
		true,
		[]byte("test"),
		map[string]int{"key": 1},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CastString(testValues[i%len(testValues)])
	}
}

// Benchmark for NowUInt function
func BenchmarkNowUInt(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NowUInt()
	}
}
