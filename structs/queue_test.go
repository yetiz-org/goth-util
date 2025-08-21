package structs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// SimpleQueue is a basic implementation of the Queue interface for testing
type SimpleQueue struct {
	items []interface{}
}

// NewSimpleQueue creates a new SimpleQueue instance
func NewSimpleQueue() *SimpleQueue {
	return &SimpleQueue{
		items: make([]interface{}, 0),
	}
}

// Push adds an element to the back of the queue
func (q *SimpleQueue) Push(v interface{}) {
	q.items = append(q.items, v)
}

// Pop removes and returns the element at the front of the queue
func (q *SimpleQueue) Pop() interface{} {
	if len(q.items) == 0 {
		return nil
	}
	
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

// Size returns the number of elements in the queue
func (q *SimpleQueue) Size() int {
	return len(q.items)
}

// IsEmpty returns true if the queue is empty
func (q *SimpleQueue) IsEmpty() bool {
	return len(q.items) == 0
}

func TestQueueInterface(t *testing.T) {
	// Test that SimpleQueue implements Queue interface
	var queue Queue = NewSimpleQueue()
	assert.NotNil(t, queue)

	// Test Push and Pop through interface
	queue.Push("hello")
	queue.Push(42)
	queue.Push(true)

	// Pop items in FIFO order
	item := queue.Pop()
	assert.Equal(t, "hello", item)

	item = queue.Pop()
	assert.Equal(t, 42, item)

	item = queue.Pop()
	assert.Equal(t, true, item)

	// Pop from empty queue should return nil
	item = queue.Pop()
	assert.Nil(t, item)
}

func TestSimpleQueueBasicOperations(t *testing.T) {
	queue := NewSimpleQueue()

	// Test initial state
	assert.True(t, queue.IsEmpty())
	assert.Equal(t, 0, queue.Size())

	// Test Push operation
	queue.Push("first")
	assert.False(t, queue.IsEmpty())
	assert.Equal(t, 1, queue.Size())

	queue.Push("second")
	assert.Equal(t, 2, queue.Size())

	// Test Pop operation (FIFO)
	item := queue.Pop()
	assert.Equal(t, "first", item)
	assert.Equal(t, 1, queue.Size())

	item = queue.Pop()
	assert.Equal(t, "second", item)
	assert.Equal(t, 0, queue.Size())
	assert.True(t, queue.IsEmpty())
}

func TestSimpleQueueWithDifferentTypes(t *testing.T) {
	queue := NewSimpleQueue()

	// Test with different data types
	queue.Push(123)
	queue.Push("string")
	queue.Push([]int{1, 2, 3})
	queue.Push(map[string]int{"key": 1})

	// Pop and verify types
	item := queue.Pop()
	assert.Equal(t, 123, item)

	item = queue.Pop()
	assert.Equal(t, "string", item)

	item = queue.Pop()
	assert.Equal(t, []int{1, 2, 3}, item)

	item = queue.Pop()
	assert.Equal(t, map[string]int{"key": 1}, item)
}

func TestSimpleQueuePopFromEmpty(t *testing.T) {
	queue := NewSimpleQueue()

	// Pop from empty queue multiple times
	for i := 0; i < 5; i++ {
		item := queue.Pop()
		assert.Nil(t, item)
		assert.True(t, queue.IsEmpty())
		assert.Equal(t, 0, queue.Size())
	}
}

func TestSimpleQueueLargeVolume(t *testing.T) {
	queue := NewSimpleQueue()
	count := 1000

	// Push many items
	for i := 0; i < count; i++ {
		queue.Push(i)
	}

	assert.Equal(t, count, queue.Size())
	assert.False(t, queue.IsEmpty())

	// Pop all items and verify order
	for i := 0; i < count; i++ {
		item := queue.Pop()
		assert.Equal(t, i, item)
	}

	assert.Equal(t, 0, queue.Size())
	assert.True(t, queue.IsEmpty())
}

func TestSimpleQueueMixedOperations(t *testing.T) {
	queue := NewSimpleQueue()

	// Mix push and pop operations
	queue.Push("a")
	queue.Push("b")

	item := queue.Pop()
	assert.Equal(t, "a", item)

	queue.Push("c")
	queue.Push("d")

	item = queue.Pop()
	assert.Equal(t, "b", item)

	item = queue.Pop()
	assert.Equal(t, "c", item)

	item = queue.Pop()
	assert.Equal(t, "d", item)

	item = queue.Pop()
	assert.Nil(t, item)
}

// GenericQueue is a generic implementation of Queue for testing with Go 1.18+ generics
type GenericQueue[T any] struct {
	items []T
}

func NewGenericQueue[T any]() *GenericQueue[T] {
	return &GenericQueue[T]{
		items: make([]T, 0),
	}
}

func (q *GenericQueue[T]) Push(v interface{}) {
	if item, ok := v.(T); ok {
		q.items = append(q.items, item)
	}
}

func (q *GenericQueue[T]) Pop() interface{} {
	if len(q.items) == 0 {
		var zero T
		return zero
	}
	
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func (q *GenericQueue[T]) Size() int {
	return len(q.items)
}

func TestGenericQueue(t *testing.T) {
	// Test with string type
	stringQueue := NewGenericQueue[string]()
	var queue Queue = stringQueue

	queue.Push("hello")
	queue.Push("world")

	item := queue.Pop()
	assert.Equal(t, "hello", item)

	item = queue.Pop()
	assert.Equal(t, "world", item)

	// Test with int type
	intQueue := NewGenericQueue[int]()
	queue = intQueue

	queue.Push(1)
	queue.Push(2)
	queue.Push(3)

	item = queue.Pop()
	assert.Equal(t, 1, item)

	item = queue.Pop()
	assert.Equal(t, 2, item)

	item = queue.Pop()
	assert.Equal(t, 3, item)
}

func TestQueueWithNilValues(t *testing.T) {
	queue := NewSimpleQueue()

	// Test pushing nil values
	queue.Push(nil)
	queue.Push("not nil")
	queue.Push(nil)

	assert.Equal(t, 3, queue.Size())

	// Pop and verify nil handling
	item := queue.Pop()
	assert.Nil(t, item)

	item = queue.Pop()
	assert.Equal(t, "not nil", item)

	item = queue.Pop()
	assert.Nil(t, item)

	// Empty queue should also return nil
	item = queue.Pop()
	assert.Nil(t, item)
}
