package concurrent

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFuture(t *testing.T) {
	fs := NewFuture(context.Background())
	for i := 0; i < 100; i++ {
		go func() {
			assert.EqualValues(t, 1, fs.Get())
			assert.EqualValues(t, true, fs.IsDone())
			assert.EqualValues(t, true, fs.IsSuccess())
		}()
	}

	go func() {
		time.Sleep(time.Microsecond)
		fs.(CompletableFuture).Complete(1)
	}()

	assert.EqualValues(t, false, fs.IsDone())
	assert.EqualValues(t, 1, fs.Get())
	assert.EqualValues(t, true, fs.IsDone())
	assert.EqualValues(t, true, fs.IsSuccess())
	assert.EqualValues(t, false, fs.IsCancelled())
	assert.EqualValues(t, false, fs.IsError())

	f := NewFuture(nil)
	go func() {
		time.Sleep(time.Microsecond)
		f.(CompletableFuture).Complete(1)
	}()

	f.AddListener(NewFutureListener(func(f Future) {
		if !f.IsSuccess() {
			assert.Fail(t, "should not go into this scope")
		}
	}))

	f.Await()

	f = NewFuture(nil)
	go func() {
		time.Sleep(time.Microsecond)
		f.(CompletableFuture).Cancel()
	}()

	f.AddListener(NewFutureListener(func(f Future) {
		if !f.IsCancelled() {
			assert.Fail(t, "should not go into this scope")
		}
	}))

	f.Await()

	f = NewFuture(nil)
	go func() {
		time.Sleep(time.Microsecond)
		f.(CompletableFuture).Fail(fmt.Errorf("fail"))
	}()

	f.AddListener(NewFutureListener(func(f Future) {
		if !f.IsError() {
			assert.Fail(t, "should not go into this scope")
		}
	}))

	assert.EqualValues(t, false, f.IsDone())
	f.Await()
	assert.EqualValues(t, "fail", f.Error().Error())

	f = NewSucceededFuture("s")
	assert.EqualValues(t, "s", f.Get())
	assert.EqualValues(t, true, f.IsDone())
	assert.EqualValues(t, true, f.IsSuccess())
	assert.EqualValues(t, nil, f.Error())
	f.AddListener(NewFutureListener(func(f Future) {
		assert.Fail(t, "should not go into this scope")
	}))

	f = NewCancelledFuture()
	assert.EqualValues(t, nil, f.Get())
	assert.EqualValues(t, true, f.IsDone())
	assert.EqualValues(t, true, f.IsCancelled())
	assert.EqualValues(t, false, f.Error() == nil)
	f.AddListener(NewFutureListener(func(f Future) {
		assert.Fail(t, "should not go into this scope")
	}))

	f = NewFailedFuture(fmt.Errorf("err"))
	assert.EqualValues(t, nil, f.Get())
	assert.EqualValues(t, true, f.IsDone())
	assert.EqualValues(t, true, f.IsError())
	assert.EqualValues(t, "err", f.Error().Error())
	f.AddListener(NewFutureListener(func(f Future) {
		assert.Fail(t, "should not go into this scope")
	}))

	ctx, cancelFunc := context.WithCancel(context.Background())
	f = NewFuture(ctx)
	f.AddListener(NewFutureListener(func(f Future) {
		assert.EqualValues(t, true, f.IsDone())
		assert.EqualValues(t, true, f.IsCancelled())
	}))

	cancelFunc()
	assert.EqualValues(t, false, f.IsDone())
	assert.EqualValues(t, false, f.IsSuccess())
	assert.EqualValues(t, false, f.IsCancelled())
	assert.EqualValues(t, false, f.IsError())
	time.Sleep(time.Millisecond)
	assert.EqualValues(t, true, f.IsDone())
	assert.EqualValues(t, true, f.IsCancelled())

	ffv := int32(0)
	f = NewFuture(nil)
	listener := NewFutureListener(func(f Future) {
		ffvP := f.Get().(*int32)
		assert.EqualValues(t, 0, *ffvP)
		atomic.StoreInt32(ffvP, 1)
	}).AddListener(NewFutureListener(func(f Future) {
		if f.IsSuccess() {
			ffvP := f.Get().(*int32)
			assert.EqualValues(t, 1, *ffvP)
			atomic.StoreInt32(ffvP, 2)
		}
	}))

	f.AddListener(listener.(FutureListener))
	f.(CompletableFuture).Complete(&ffv)
	assert.EqualValues(t, 2, ffv)
}
