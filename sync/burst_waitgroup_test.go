package sync

import (
	"testing"
	"time"
)

func TestBurstWaitGroup(t *testing.T) {
	b := &BurstWaitGroup{}
	b.Add(100)
	go func() {
		time.Sleep(time.Millisecond)
		for i := 0; i < 100; i++ {
			go func() { b.Done() }()
		}

		b.Burst()
	}()

	b.Wait()
}
