package sync

import "sync"

type Queue struct {
	head *queueElement
	tail *queueElement
	len  int
	op   sync.Mutex
}

func (q *Queue) Push(obj interface{}) *Queue {
	if obj == nil {
		return q
	}

	q.op.Lock()
	defer q.op.Unlock()
	if q.head == nil {
		q.head = &queueElement{
			prev: nil,
			next: nil,
			obj:  obj,
		}

		q.tail = q.head
	} else {
		ne := &queueElement{
			prev: q.tail,
			next: nil,
			obj:  obj,
		}

		q.tail.next = ne
		q.tail = ne
	}

	q.len++
	return q
}

func (q *Queue) Pop() interface{} {
	q.op.Lock()
	defer q.op.Unlock()
	if q.head == nil {
		return nil
	}

	rtn := q.head
	if q.head.next == nil {
		q.head = nil
	} else {
		q.head = q.head.next
		q.head.prev = nil
	}

	rtn.next = nil
	q.len--
	return rtn.obj
}

func (q *Queue) Reset() *Queue {
	q.op.Lock()
	defer q.op.Unlock()
	q.head = nil
	q.tail = nil
	q.len = 0
	return q
}

func (q *Queue) Len() int {
	return q.len
}

type queueElement struct {
	prev, next *queueElement
	obj        interface{}
}
