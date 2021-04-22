package sync

import (
	"testing"

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
}
