package concurrent

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

const (
	stateWait = iota
	stateSuccess
	stateCancel
	stateFail
)

type Future interface {
	Get() interface{}
	Sync() Future
	Await() Future
	IsDone() bool
	IsSuccess() bool
	IsCancelled() bool
	IsError() bool
	Error() error
	AddListener(listener FutureListener) Future
}

type CompletableFuture interface {
	Complete(obj interface{})
	Cancel()
	Fail(err error)
	callListener()
}

func NewFuture(ctx context.Context) Future {
	f := &DefaultFuture{}
	if ctx == nil {
		ctx = context.Background()
	}

	f.ctx, f.cancel = context.WithCancel(ctx)
	return f
}

func NewSuccessedFuture(obj interface{}) Future {
	f := &DefaultFuture{}
	f.ctx, f.cancel = context.WithCancel(context.Background())
	f.Complete(obj)
	return f
}

func NewCancelledFuture() Future {
	f := &DefaultFuture{}
	f.ctx, f.cancel = context.WithCancel(context.Background())
	f.Cancel()
	return f
}

func NewFailedFuture() Future {
	f := &DefaultFuture{}
	f.ctx, f.cancel = context.WithCancel(context.Background())
	f.Fail(fmt.Errorf("empty future"))
	return f
}

type DefaultFuture struct {
	obj       interface{}
	state     int32
	err       error
	ctx       context.Context
	cancel    context.CancelFunc
	opL       sync.Mutex
	listeners []FutureListener
	once      sync.Once
}

func (f *DefaultFuture) Get() interface{} {
	if f.IsDone() {
		return f.obj
	}

	<-f.ctx.Done()
	if !f.IsDone() {
		atomic.StoreInt32(&f.state, stateCancel)
		f.err = f.ctx.Err()
		f.callListener()
	}

	return f.obj
}

func (f *DefaultFuture) Sync() Future {
	f.Get()
	return f
}

func (f *DefaultFuture) Await() Future {
	return f.Sync()
}

func (f *DefaultFuture) IsDone() bool {
	return atomic.LoadInt32(&f.state) != stateWait
}

func (f *DefaultFuture) IsSuccess() bool {
	return atomic.LoadInt32(&f.state) == stateSuccess
}

func (f *DefaultFuture) IsCancelled() bool {
	return atomic.LoadInt32(&f.state) == stateCancel
}

func (f *DefaultFuture) IsError() bool {
	return atomic.LoadInt32(&f.state) == stateFail
}

func (f *DefaultFuture) Error() error {
	return f.err
}

func (f *DefaultFuture) AddListener(listener FutureListener) Future {
	f.opL.Lock()
	defer f.opL.Unlock()
	f.listeners = append(f.listeners, listener)
	return f
}

func (f *DefaultFuture) Complete(obj interface{}) {
	f.opL.Lock()
	defer f.opL.Unlock()
	if !f.IsDone() {
		atomic.StoreInt32(&f.state, stateSuccess)
		f.obj = obj
		f.callListener()
		f.cancel()
	}
}

func (f *DefaultFuture) Cancel() {
	f.opL.Lock()
	defer f.opL.Unlock()
	if !f.IsDone() {
		f.cancel()
	}
}

func (f *DefaultFuture) Fail(err error) {
	f.opL.Lock()
	defer f.opL.Unlock()
	if !f.IsDone() {
		atomic.StoreInt32(&f.state, stateFail)
		f.err = err
		f.callListener()
		f.cancel()
	}
}

func (f *DefaultFuture) callListener() {
	f.once.Do(func() {
		for _, listener := range f.listeners {
			listener.OperationCompleted(f)
		}
	})
}

type FutureListener interface {
	Future
	OperationCompleted(f Future)
}

type _FutureListener struct {
	DefaultFuture
	f func(f Future)
}

func (l *_FutureListener) OperationCompleted(f Future) {
	l.f(f)
}

func NewFutureListener(f func(f Future)) FutureListener {
	return &_FutureListener{
		f: f,
	}
}
