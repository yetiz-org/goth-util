package sync

import (
	"math/rand"
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
	for i := 0; i < 32000; i++ {
		go func() {
			b.Add(1)
			time.Sleep(time.Microsecond * time.Duration(rand.Int()%1000))
			b.Done()
		}()
	}

	b.Wait()
}
