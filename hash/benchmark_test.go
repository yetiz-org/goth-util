package hash

import (
	"testing"
	"time"
)

// BenchmarkTimeHash tests the performance of TimeHash function
func BenchmarkTimeHash(b *testing.B) {
	data := []byte("Hello, World! This is a test data for time hash benchmarking.")
	timestamp := time.Now().Unix()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = TimeHash(data, timestamp)
	}
}

// BenchmarkTimeHashSmallData tests TimeHash performance with small data
func BenchmarkTimeHashSmallData(b *testing.B) {
	data := []byte("small")
	timestamp := time.Now().Unix()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = TimeHash(data, timestamp)
	}
}

// BenchmarkTimeHashLargeData tests TimeHash performance with larger data
func BenchmarkTimeHashLargeData(b *testing.B) {
	data := make([]byte, 1024) // 1KB data
	for i := range data {
		data[i] = byte(i % 256)
	}
	timestamp := time.Now().Unix()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = TimeHash(data, timestamp)
	}
}

// BenchmarkCryptoTimeHash tests the performance of CryptoTimeHash function
func BenchmarkCryptoTimeHash(b *testing.B) {
	data := []byte("Hello, World! This is a test data for crypto time hash benchmarking.")
	timestamp := time.Now().Unix()
	key := []byte("test-encryption-key-for-benchmarking")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CryptoTimeHash(data, timestamp, key)
	}
}

// BenchmarkCryptoTimeHashSmallData tests CryptoTimeHash performance with small data
func BenchmarkCryptoTimeHashSmallData(b *testing.B) {
	data := []byte("small")
	timestamp := time.Now().Unix()
	key := []byte("test-key")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CryptoTimeHash(data, timestamp, key)
	}
}

// BenchmarkEncrypt tests the performance of _Encrypt function
func BenchmarkEncrypt(b *testing.B) {
	data := []byte("Hello, World! This is a test data for encryption benchmarking.")
	key := []byte("test-encryption-key-for-benchmarking")
	
	// Pad data to AES block size (16 bytes)
	blockSize := 16
	padLen := blockSize - (len(data) % blockSize)
	if padLen != blockSize {
		padding := make([]byte, padLen)
		data = append(data, padding...)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = _Encrypt(key, data)
	}
}

// BenchmarkDecrypt tests the performance of _Decrypt function
func BenchmarkDecrypt(b *testing.B) {
	data := []byte("Hello, World! This is a test data for decryption benchmarking.")
	key := []byte("test-encryption-key-for-benchmarking")
	
	// Pad data to AES block size (16 bytes)
	blockSize := 16
	padLen := blockSize - (len(data) % blockSize)
	if padLen != blockSize {
		padding := make([]byte, padLen)
		data = append(data, padding...)
	}
	
	encrypted := _Encrypt(key, data)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = _Decrypt(key, encrypted)
	}
}

// BenchmarkValidateTimeHash tests the performance of ValidateTimeHash function
func BenchmarkValidateTimeHash(b *testing.B) {
	data := []byte("test data for validation")
	timestamp := time.Now().Unix()
	encoded := TimeHash(data, timestamp)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateTimeHash(encoded)
	}
}

// BenchmarkDataOfTimeHash tests the performance of DataOfTimeHash function
func BenchmarkDataOfTimeHash(b *testing.B) {
	data := []byte("test data for extraction")
	timestamp := time.Now().Unix()
	encoded := TimeHash(data, timestamp)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DataOfTimeHash(encoded)
	}
}

// BenchmarkTimestampOfTimeHash tests the performance of TimestampOfTimeHash function
func BenchmarkTimestampOfTimeHash(b *testing.B) {
	data := []byte("test data for timestamp extraction")
	timestamp := time.Now().Unix()
	encoded := TimeHash(data, timestamp)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = TimestampOfTimeHash(encoded)
	}
}
