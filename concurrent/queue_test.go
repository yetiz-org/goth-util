package concurrent

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	q := &Queue{}
	q.Push(1)
	q.Push(2)
	assert.EqualValues(t, 2, q.Len())
	assert.EqualValues(t, q.Pop(), 1)
	assert.EqualValues(t, q.Pop(), 2)
	assert.EqualValues(t, 0, q.Len())
	q.Push(1)
	q.Push(2)
	assert.EqualValues(t, q.Pop(), 1)
	q.Reset()
	assert.EqualValues(t, 0, q.Len())
	assert.Nil(t, q.Pop())
}

func TestBlockingQueue(t *testing.T) {
	q := &BlockingQueue{}
	q.Push(1)
	q.Push(2)
	assert.EqualValues(t, 2, q.Len())
	assert.EqualValues(t, q.Pop(), 1)
	assert.EqualValues(t, q.Pop(), 2)
	assert.EqualValues(t, 0, q.Len())
	q.Push(1)
	q.Push(2)
	assert.EqualValues(t, q.Pop(), 1)
	q.Reset()
	assert.EqualValues(t, 0, q.Len())
	go func() {
		time.Sleep(time.Millisecond)
		q.Push(1)
	}()

	assert.NotNil(t, q.Pop())
	assert.EqualValues(t, 0, q.Len())
}
