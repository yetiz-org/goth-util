package concurrent

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMutex(t *testing.T) {
	m := Mutex{}
	ts := time.Now()
	assert.True(t, m.TryLock())
	assert.False(t, m.TryLock())
	m.Unlock()
	assert.True(t, m.TryLock())
	assert.False(t, m.TryLock())
	m.Unlock()
	assert.True(t, m.TryLock())
	m.Unlock()
	assert.True(t, m.TryLock())
	assert.False(t, m.TryLockTimeout(time.Millisecond))
	m.Unlock()
	ts = time.Now()
	assert.True(t, m.TryLockTimeout(time.Millisecond))
	assert.True(t, time.Now().Sub(ts) > time.Microsecond)
	m.Unlock()
	for i := 0; i < 100; i++ {
		go func() {
			assert.Panics(t, func() {
				m.Unlock()
			})
		}()
	}

	time.Sleep(time.Second)
	m.Lock()
	assert.False(t, m.TryLock())
	assert.False(t, m.TryLockTimeout(time.Millisecond))
	state := int32(0)
	for i := 0; i < 100; i++ {
		go func(state *int32) {
			time.Sleep(time.Millisecond)
			if m.TryLock() {
				atomic.AddInt32(state, 1)
			}
		}(&state)
	}

	time.Sleep(time.Millisecond)
	m.Unlock()
	time.Sleep(time.Millisecond)
	assert.False(t, m.TryLock())
	assert.Equal(t, int32(1), state)
	m.Unlock()
	for i := 0; i < 100; i++ {
		go func(state *int32) {
			if m.TryLock() {
				atomic.AddInt32(state, 1)
			}
		}(&state)
	}

	assert.Equal(t, int32(2), state)
}
