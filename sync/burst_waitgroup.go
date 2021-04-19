package sync

import (
	"sync"
)

type BurstWaitGroup struct {
	wg    sync.WaitGroup
	delta int
	l     sync.Mutex
}

func (w *BurstWaitGroup) Add(delta int) {
	if w.delta+delta < 0 {
		return
	}

	w.l.Lock()
	defer w.l.Unlock()
	if w.delta+delta < 0 {
		return
	}

	w.delta += delta
	w.wg.Add(delta)
}

func (w *BurstWaitGroup) Done() {
	w.Add(-1)
}

func (w *BurstWaitGroup) Wait() {
	w.wg.Wait()
}

func (w *BurstWaitGroup) Remain() int {
	return w.delta
}

func (w *BurstWaitGroup) Burst() {
	w.Add(w.delta * -1)

	if w.delta > 0 {
		for i := 0; i < w.delta; i++ {
			w.Add(-1)
		}
	}
}
