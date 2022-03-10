package hash

import (
	"crypto/aes"
	"crypto/cipher"
	rand2 "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"io"
	"math/rand"

	"github.com/kklab-com/goth-base62"
	buf "github.com/kklab-com/goth-bytebuf"
)

var TimeHashBase = []byte{75, 79, 78, 83, 73, 84, 69, 89}
var CryptoTimeHashPadding = byte(0x59)

/*
encode data with timestamp
`data` can't be nil or empty slice,
`timestamp` can't less than or equal to 0
*/
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
				bs[i] = byte(rand.Int31n(255))
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
	r[0] = 0x59
	r[1] = v
	r[2] = pad
	r[rl-1] = 0x53
	for i := 0; i < 8; i++ {
		for j := 0; j < align; j++ {
			r[3+i*(align+1)+j] = data[i*align+j] ^ bs[i]
		}

		r[2+(i+1)*(align+1)] = bs[i]
	}

	for i := 1; i < rl-1; i++ {
		r[rl-1] ^= r[i]
	}

	r[0] = byte((int(r[0]) + int(r[rl-1])) % 255)
	r[1] ^= r[0]
	r[2] ^= r[0]
	return base62.ShiftEncoding.EncodeToString(r)
}

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

		bb := buf.NewByteBuf([]byte{c})
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
				bs[i] = byte(rand.Int31n(255))
			}
		}

		data = bs
	}

	rl := dpl + 12
	r := make([]byte, rl)
	r[0] = 0x59
	r[1] = v
	r[2] = pad
	r[rl-1] = 0x53
	for i := 0; i < 8; i++ {
		for j := 0; j < align; j++ {
			r[3+i*(align+1)+j] = data[i*align+j] ^ tbs[i]
		}

		r[2+(i+1)*(align+1)] = tbs[i]
	}

	for i := 1; i < rl-1; i++ {
		r[rl-1] ^= r[i]
	}

	r[0] = byte((int(r[0]) + int(r[rl-1])) % 255)
	r[1] ^= r[0]
	r[2] ^= r[0]
	return base62.ShiftEncoding.EncodeToString(r)
}

func _Encrypt(key []byte, data []byte) []byte {
	if data == nil {
		return nil
	}

	key = func() []byte {
		r := make([]byte, 32)
		for i, k := 0, sha256.Sum256(key); i < 32; i++ {
			r[i] = k[i]
		}

		return r
	}()

	block, _ := aes.NewCipher(key)
	rtn := make([]byte, aes.BlockSize+len(data))
	iv := rtn[:aes.BlockSize]
	if _, err := io.ReadFull(rand2.Reader, iv); err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(rtn[aes.BlockSize:], data)
	return rtn
}

func _Decrypt(key []byte, encrypted []byte) []byte {
	if len(encrypted) < aes.BlockSize || len(encrypted)%aes.BlockSize != 0 {
		return nil
	}

	key = func() []byte {
		r := make([]byte, 32)
		for i, k := 0, sha256.Sum256(key); i < 32; i++ {
			r[i] = k[i]
		}

		return r
	}()

	block, _ := aes.NewCipher(key)
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

	var c byte = 0x53
	for i := 1; i < dl-1; i++ {
		c ^= d[i]
	}

	if c != d[dl-1] {
		return false
	}

	if byte((int(c)+0x59)%255) != d[0] {
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
