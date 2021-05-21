package sync

import "sync"

var queuePool = &sync.Pool{
	New: func() interface{} { return &queueElement{} },
}

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
		qe := queuePool.Get().(*queueElement)
		qe.prev = nil
		qe.next = nil
		qe.obj = obj
		q.head = qe
		q.tail = q.head
	} else {
		qe := queuePool.Get().(*queueElement)
		qe.prev = q.tail
		qe.next = nil
		qe.obj = obj
		q.tail.next = qe
		q.tail = qe
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

	obj := rtn.obj
	rtn.prev = nil
	rtn.next = nil
	rtn.obj = nil
	queuePool.Put(rtn)
	q.len--
	return obj
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
