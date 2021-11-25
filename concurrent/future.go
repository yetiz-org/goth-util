package concurrent

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync/atomic"
	"time"
)

const (
	stateWait = iota
	stateSuccess
	stateCancel
	stateFail
)

type Future interface {
	Get() interface{}
	GetTimeout(timeout time.Duration) interface{}
	Done() <-chan struct{}
	Await() Future
	AwaitTimeout(timeout time.Duration) Future
	IsDone() bool
	IsSuccess() bool
	IsCancelled() bool
	IsError() bool
	Error() error
	Ctx() context.Context
	AddListener(listener FutureListener) Future
	Completable() CompletableFuture
	Immutable() ImmutableFuture
}

type ImmutableFuture interface {
	Get() interface{}
	GetTimeout(timeout time.Duration) interface{}
	Done() <-chan struct{}
	Await() Future
	AwaitTimeout(timeout time.Duration) Future
	IsDone() bool
	IsSuccess() bool
	IsCancelled() bool
	IsError() bool
	Error() error
	Ctx() context.Context
	AddListener(listener FutureListener) Future
}

type CompletableFuture interface {
	Complete(obj interface{})
	Cancel()
	Fail(err error)
	callListener()
}

type CarrierFuture interface {
	Payload() interface{}
}

func NewCarrierFuture(obj interface{}) Future {
	f := NewFuture(nil)
	f.(*DefaultFuture).obj = obj
	return f
}

func NewFuture(ctx context.Context) Future {
	var f = &DefaultFuture{
		listeners: []FutureListener{},
	}

	if ctx == nil {
		f.ctx, f.cancel = context.WithCancel(context.Background())
	} else {
		f.ctx, f.cancel = context.WithCancel(ctx)
		go func(f *DefaultFuture) {
			f._waitCancelJudge(0)
		}(f)
	}

	return f
}

func NewSucceededFuture(obj interface{}) Future {
	f := NewFuture(nil)
	f.(CompletableFuture).Complete(obj)
	return f
}

func NewCancelledFuture() Future {
	f := NewFuture(nil)
	f.(CompletableFuture).Cancel()
	return f
}

func NewFailedFuture(err error) Future {
	f := NewFuture(nil)
	f.(CompletableFuture).Fail(err)
	return f
}

type DefaultFuture struct {
	obj       interface{}
	state     int32
	err       error
	ctx       context.Context
	cancel    context.CancelFunc
	listeners []FutureListener
}

func (f *DefaultFuture) _waitCancelJudge(timeout time.Duration) (done bool) {
	return f._waitJudge(timeout, true)
}

func (f *DefaultFuture) _waitJudge(timeout time.Duration, withCancel bool) (done bool) {
	if timeout == 0 {
		<-f.ctx.Done()
		if withCancel {
			f._cancelJudge()
		}

		return true
	} else {
		select {
		case <-f.ctx.Done():
			if withCancel {
				f._cancelJudge()
			}

			return true
		case <-time.After(timeout):
			return false
		}
	}
}

func (f *DefaultFuture) _cancelJudge() {
	if err := f.ctx.Err(); err != nil {
		f.Fail(err)
	} else {
		f.Cancel()
	}
}

func (f *DefaultFuture) Get() interface{} {
	return f.GetTimeout(0)
}

func (f *DefaultFuture) GetTimeout(timeout time.Duration) interface{} {
	if f.IsDone() {
		return f.obj
	}

	if f._waitCancelJudge(timeout) {
		return f.obj
	}

	return nil
}

func (f *DefaultFuture) Done() <-chan struct{} {
	return f.ctx.Done()
}

func (f *DefaultFuture) Await() Future {
	f.Get()
	return f
}

func (f *DefaultFuture) AwaitTimeout(timeout time.Duration) Future {
	f._waitJudge(timeout, false)
	return f
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

func (f *DefaultFuture) Ctx() context.Context {
	return f.ctx
}

func (f *DefaultFuture) AddListener(listener FutureListener) Future {
	if listener == nil {
		return f
	}

	f.listeners = append(f.listeners, listener)
	return f
}

func (f *DefaultFuture) self() Future {
	return f
}

func (f *DefaultFuture) Completable() CompletableFuture {
	return f.self().(CompletableFuture)
}

func (f *DefaultFuture) Immutable() ImmutableFuture {
	return f.self().(ImmutableFuture)
}

func (f *DefaultFuture) Complete(obj interface{}) {
	if atomic.CompareAndSwapInt32(&f.state, stateWait, stateSuccess) {
		if obj != nil {
			f.obj = obj
		}

		f.callListener()
		f.cancel()
	}
}

func (f *DefaultFuture) Cancel() {
	if atomic.CompareAndSwapInt32(&f.state, stateWait, stateCancel) {
		f.callListener()
		f.cancel()
	}
}

func (f *DefaultFuture) Fail(err error) {
	if atomic.CompareAndSwapInt32(&f.state, stateWait, stateFail) {
		f.err = err
		f.callListener()
		f.cancel()
	}
}

func (f *DefaultFuture) callListener() {
	if f.listeners == nil {
		f.listeners = []FutureListener{}
	}

	for _, listener := range f.listeners {
		if listener == nil {
			continue
		}

		listener.OperationCompleted(f)
	}
}

func (f *DefaultFuture) Payload() interface{} {
	return f.obj
}

type FutureListener interface {
	Future
	OperationCompleted(f Future)
}

type _FutureListener struct {
	Future
	f func(f Future)
}

func (l *_FutureListener) OperationCompleted(f Future) {
	defer func() {
		if v := recover(); v != nil {
			println(v)
			println(string(debug.Stack()))
		}
	}()

	if l.f == nil {
		l.Future.(CompletableFuture).Fail(fmt.Errorf("nil f in future listener"))
		return
	}

	func(f Future) {
		defer func() {
			if v := recover(); v != nil {
				switch cast := v.(type) {
				case error:
					l.Future.(CompletableFuture).Fail(cast)
				default:
					l.Future.(CompletableFuture).Fail(fmt.Errorf("%v", cast))
				}
			}
		}()

		l.f(f)
		l.Future.(CompletableFuture).Complete(f.Get())
	}(f)
}

func (l *_FutureListener) AddListener(listener FutureListener) Future {
	if listener == nil {
		return l
	}

	ll := listener
	l.Future.(*DefaultFuture).listeners = append(l.Future.(*DefaultFuture).listeners, ll)
	return l
}

func NewFutureListener(f func(f Future)) FutureListener {
	lf := f
	return &_FutureListener{
		Future: NewFuture(nil),
		f:      lf,
	}
}
