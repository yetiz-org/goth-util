// Package hash provides time-based hashing utilities with optional encryption support.
// It implements a custom encoding scheme that embeds timestamp information along with data,
// allowing for data integrity verification and timestamp extraction.
package hash

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	rand2 "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"io"
	"math/rand"

	"github.com/yetiz-org/goth-base62"
)

// TimeHashBase is the base key used for XOR operations in the time hash algorithm.
// This constant provides a fixed seed for consistent encoding/decoding operations.
var TimeHashBase = []byte{75, 79, 78, 83, 73, 84, 69, 89}

// CryptoTimeHashPadding is the padding byte used in encrypted time hash operations.
// This value is used for both padding data and as part of the integrity check.
var CryptoTimeHashPadding = byte(0x59)

// CryptoTimeHashXBit is the XOR bit used for checksum calculations in time hashes.
// This value helps ensure data integrity during encoding and decoding operations.
var CryptoTimeHashXBit = byte(0x53)

// TimeHash encodes data with an embedded timestamp using a custom algorithm.
// The function combines the data with the timestamp and applies XOR operations
// with padding and checksum validation to create a secure, time-stamped hash.
//
// Parameters:
//   - data: The byte slice to be encoded (cannot be nil or empty)
//   - timestamp: Unix timestamp in seconds (must be greater than 0)
//
// Returns:
//   - A base62-encoded string containing the hashed data with embedded timestamp
//   - Empty string if input validation fails
//
// The encoding process:
//  1. Validates input parameters
//  2. Applies padding to align data to 8-byte boundaries
//  3. XORs data segments with timestamp-derived keys
//  4. Adds integrity checksums
//  5. Encodes result using base62 encoding
func TimeHash(data []byte, timestamp int64) string {
	if data == nil || len(data) == 0 || timestamp <= 0 {
		return ""
	}

	//fck, v, pad, data with pad split with timestamp, bck
	var v byte = 0x01
	var pad = byte((8 - (len(data) % 8)) % 8)
	var dpl = len(data) + int(pad)
	var align = dpl / 8

	if dpl > len(data) {
		bs := make([]byte, dpl)
		for i, dl := 0, len(data); i < dpl; i++ {
			if i < dl {
				bs[i] = data[i]
			} else {
				bs[i] = byte(rand.Int31n(256))
			}
		}

		data = bs
	}

	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, uint64(timestamp))
	for i, b := range bs {
		bs[i] = b ^ TimeHashBase[i]
	}

	rl := dpl + 12
	r := make([]byte, rl)
	r[0] = CryptoTimeHashPadding
	r[1] = v
	r[2] = pad
	r[rl-1] = CryptoTimeHashXBit
	for i := 0; i < 8; i++ {
		for j := 0; j < align; j++ {
			r[3+i*(align+1)+j] = data[i*align+j] ^ bs[i]
		}

		r[2+(i+1)*(align+1)] = bs[i]
	}

	for i := 1; i < rl-1; i++ {
		r[rl-1] ^= r[i]
	}

	r[0] = byte((int(r[0]) + int(r[rl-1])) % 256)
	r[1] ^= r[0]
	r[2] ^= r[0]
	return base62.ShiftEncoding.EncodeToString(r)
}

// CryptoTimeHash encodes data with an embedded timestamp using AES encryption.
// This function provides enhanced security by encrypting the data before applying
// the time hash algorithm, making it suitable for sensitive information.
//
// Parameters:
//   - data: The byte slice to be encoded and encrypted (cannot be nil or empty)
//   - timestamp: Unix timestamp in seconds (must be greater than 0)
//   - key: Encryption key (cannot be nil or empty, will be SHA256 hashed to 32-byte key)
//
// Returns:
//   - A base62-encoded string containing the encrypted and hashed data with embedded timestamp
//   - Empty string if input validation fails or encryption fails
//
// The encoding process:
//  1. Validates all input parameters
//  2. Encrypts the data using AES-256-CBC with the provided key
//  3. Applies the TimeHash algorithm to the encrypted data
//  4. Returns the base62-encoded result
//
// Security features:
//   - AES-256-CBC encryption with random IV
//   - SHA256 key derivation for consistent 32-byte keys
//   - Integrity validation through checksums
func CryptoTimeHash(data []byte, timestamp int64, key []byte) string {
	if data == nil || len(data) == 0 || timestamp <= 0 || key == nil || len(key) == 0 {
		return ""
	}

	tbs := make([]byte, 8)
	binary.LittleEndian.PutUint64(tbs, uint64(timestamp))
	for i, b := range tbs {
		tbs[i] = b ^ TimeHashBase[i]
	}

	//fck, v, pad, data with pad split with timestamp, bck
	var v byte = 0x02
	data = func(tbs []byte) []byte {
		c := byte(0x00)
		for _, b := range tbs {
			c ^= b
		}

		bb := bytes.NewBuffer([]byte{c})
		if dr := (16 - ((len(data) + 2) % 16)) % 16; dr > 0 {
			bb.WriteByte(byte(dr))
			bb.Write(data)
			for i := 0; i < dr; i++ {
				bb.WriteByte(CryptoTimeHashPadding)
			}
		} else {
			bb.WriteByte(0x00)
			bb.Write(data)
		}

		return bb.Bytes()
	}(tbs)

	if d := _Encrypt(key, data); d != nil {
		data = d
	} else {
		return ""
	}

	var pad = byte((8 - (len(data) % 8)) % 8)
	var dpl = len(data) + int(pad)
	var align = dpl / 8

	if dpl > len(data) {
		bs := make([]byte, dpl)
		for i, dl := 0, len(data); i < dpl; i++ {
			if i < dl {
				bs[i] = data[i]
			} else {
				bs[i] = byte(rand.Int31n(256))
			}
		}

		data = bs
	}

	rl := dpl + 12
	r := make([]byte, rl)
	r[0] = CryptoTimeHashPadding
	r[1] = v
	r[2] = pad
	r[rl-1] = CryptoTimeHashXBit
	for i := 0; i < 8; i++ {
		for j := 0; j < align; j++ {
			r[3+i*(align+1)+j] = data[i*align+j] ^ tbs[i]
		}

		r[2+(i+1)*(align+1)] = tbs[i]
	}

	for i := 1; i < rl-1; i++ {
		r[rl-1] ^= r[i]
	}

	r[0] = byte((int(r[0]) + int(r[rl-1])) % 256)
	r[1] ^= r[0]
	r[2] ^= r[0]
	return base62.ShiftEncoding.EncodeToString(r)
}

// _Encrypt is an internal function that encrypts data using AES-256-CBC encryption.
// The key is automatically hashed using SHA256 to ensure a consistent 32-byte key length.
// A random IV is generated for each encryption operation to ensure security.
// Optimized version with reduced memory allocations.
//
// Parameters:
//   - key: The encryption key (will be SHA256 hashed)
//   - data: The data to encrypt
//
// Returns:
//   - Encrypted data with IV prepended, or nil if encryption fails
func _Encrypt(key []byte, data []byte) []byte {
	if data == nil {
		return nil
	}

	// Optimize key derivation - avoid extra allocation
	keyHash := sha256.Sum256(key)
	
	block, _ := aes.NewCipher(keyHash[:])
	rtn := make([]byte, aes.BlockSize+len(data))
	iv := rtn[:aes.BlockSize]
	if _, err := io.ReadFull(rand2.Reader, iv); err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(rtn[aes.BlockSize:], data)
	return rtn
}

// _Decrypt is an internal function that decrypts AES-256-CBC encrypted data.
// The key is automatically hashed using SHA256 to ensure a consistent 32-byte key length.
// The function expects the IV to be prepended to the encrypted data.
// Optimized version with reduced memory allocations.
//
// Parameters:
//   - key: The decryption key (will be SHA256 hashed)
//   - encrypted: The encrypted data with IV prepended
//
// Returns:
//   - Decrypted data, or nil if decryption fails or input is invalid
func _Decrypt(key []byte, encrypted []byte) []byte {
	if len(encrypted) < aes.BlockSize || len(encrypted)%aes.BlockSize != 0 {
		return nil
	}

	// Optimize key derivation - avoid extra allocation
	keyHash := sha256.Sum256(key)
	
	block, _ := aes.NewCipher(keyHash[:])
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encrypted, encrypted)
	return encrypted
}

func DataOfTimeHash(encoded string) []byte {
	if !ValidateTimeHash(encoded) {
		return nil
	}

	d := base62.ShiftEncoding.DecodeString(encoded)
	var dpl = len(d) - 12
	var align = dpl / 8
	var pad = int(d[2] ^ d[0])

	dp := make([]byte, dpl)
	for i := 0; i < 8; i++ {
		for j := 0; j < align; j++ {
			dp[i*align+j] = d[3+i*(align+1)+j] ^ d[2+(i+1)*(align+1)]
		}
	}

	return dp[:len(dp)-pad]
}

func TimestampOfTimeHash(encoded string) int64 {
	if !ValidateTimeHash(encoded) {
		return 0
	}

	d := base62.ShiftEncoding.DecodeString(encoded)
	var align = (len(d) - 12) / 8
	bs := make([]byte, 8)

	for i := range bs {
		bs[i] = d[2+(i+1)*(align+1)] ^ TimeHashBase[i]
	}

	v := int64(binary.LittleEndian.Uint64(bs))
	return v
}

func ValidateTimeHash(encoded string) bool {
	if encoded == "" {
		return false
	}

	d := base62.ShiftEncoding.DecodeString(encoded)
	dl := len(d)
	if dl < 20 {
		return false
	}

	var c byte = CryptoTimeHashXBit
	for i := 1; i < dl-1; i++ {
		c ^= d[i]
	}

	if c != d[dl-1] {
		return false
	}

	if byte((int(c)+int(CryptoTimeHashPadding))%256) != d[0] {
		return false
	}

	return true
}

func DataOfCryptoTimeHash(encoded string, key []byte) []byte {
	if !ValidateTimeHash(encoded) {
		return nil
	}

	d := base62.ShiftEncoding.DecodeString(encoded)
	if d[0]^d[1] != 0x02 {
		return nil
	}

	var dpl = len(d) - 12
	var align = dpl / 8
	var pad = int(d[2] ^ d[0])

	dp := make([]byte, dpl)
	for i := 0; i < 8; i++ {
		for j := 0; j < align; j++ {
			dp[i*align+j] = d[3+i*(align+1)+j] ^ d[2+(i+1)*(align+1)]
		}
	}

	tbs := make([]byte, 8)
	for i := range tbs {
		tbs[i] = d[2+(i+1)*(align+1)]
	}

	tbc := func() byte {
		c := byte(0x00)
		for _, b := range tbs {
			c ^= b
		}

		return c
	}()

	data := _Decrypt(key, dp[:len(dp)-pad])
	if tbc != data[0] {
		return nil
	}

	if data[1] == 0x00 {
		return data[2:]
	}

	return data[2 : 2+len(data[2:])-int(data[1])]
}

func FindDataOfTimeHash(encoded string, key []byte) []byte {
	if !ValidateTimeHash(encoded) {
		return nil
	}

	d := base62.ShiftEncoding.DecodeString(encoded)
	switch d[0] ^ d[1] {
	case 0x01:
		return DataOfTimeHash(encoded)
	case 0x02:
		return DataOfCryptoTimeHash(encoded, key)
	}

	return nil
}
