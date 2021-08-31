package concurrent

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

func (q *Queue) Push(obj interface{}) {
	if obj == nil {
		return
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

type BlockingQueue struct {
	Queue
	bwg BurstWaitGroup
}

func (q *BlockingQueue) Push(obj interface{}) {
	if obj == nil {
		return
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
	if q.bwg.Remain() > 0 {
		q.bwg.Burst()
	}
}

func (q *BlockingQueue) Pop() interface{} {
	q.op.Lock()
	if q.head == nil {
		q.bwg.Add(1)
		q.op.Unlock()
		q.bwg.Wait()
		return q.Pop()
	}

	defer q.op.Unlock()
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

func (q *BlockingQueue) Reset() *BlockingQueue {
	q.op.Lock()
	defer q.op.Unlock()
	q.head = nil
	q.tail = nil
	q.len = 0
	q.bwg.Burst()
	return q
}
