// Package structs provides common data structure interfaces and implementations.
package structs

// Queue defines a basic queue (FIFO - First In, First Out) interface.
// This interface provides the fundamental operations for queue data structures.
type Queue interface {
	// Pop removes and returns the element at the front of the queue.
	// Returns nil if the queue is empty.
	Pop() interface{}
	
	// Push adds an element to the back of the queue.
	// The element v can be of any type.
	Push(v interface{})
}
