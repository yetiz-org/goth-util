package hash

import (
	"crypto/rand"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeHash(t *testing.T) {
	for i := 1; i < 1024; i++ {
		bs := make([]byte, i)
		io.ReadFull(rand.Reader, bs)
		ts := time.Now().Unix()
		s := TimeHash(bs, ts)

		assert.EqualValues(t, ts, TimestampOfTimeHash(s))
		assert.EqualValues(t, bs, DataOfTimeHash(s))
		assert.EqualValues(t, bs, FindDataOfTimeHash(s, nil))
	}
}

func TestCryptoTimeHash(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	assert.Equal(t, key, _Decrypt(key, _Encrypt(key, key)))
	for i := 1; i < 1024; i++ {
		bs := make([]byte, i)
		io.ReadFull(rand.Reader, bs)
		ts := time.Now().Unix()
		s := CryptoTimeHash(bs, ts, key)

		assert.EqualValues(t, ts, TimestampOfTimeHash(s))
		assert.EqualValues(t, bs, DataOfCryptoTimeHash(s, key))
		assert.EqualValues(t, bs, FindDataOfTimeHash(s, key))
	}
}
