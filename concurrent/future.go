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
	Done() <-chan struct{}
	Await() Future
	IsDone() bool
	IsSuccess() bool
	IsCancelled() bool
	IsError() bool
	Error() error
	Ctx() context.Context
	AddListener(listener FutureListener) Future
	Completable() CompletableFuture
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
	var f = &DefaultFuture{}
	if ctx == nil {
		f.ctx, f.cancel = context.WithCancel(context.Background())
	} else {
		f.ctx, f.cancel = context.WithCancel(ctx)
		go func(f *DefaultFuture) {
			f._waitCancelJudge()
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
	opL       sync.Mutex
	listeners []FutureListener
	once      sync.Once
}

func (f *DefaultFuture) _waitCancelJudge() {
	<-f.ctx.Done()
	if atomic.CompareAndSwapInt32(&f.state, stateWait, stateCancel) {
		f.err = f.ctx.Err()
		f.callListener()
	}
}

func (f *DefaultFuture) Get() interface{} {
	if f.IsDone() {
		return f.obj
	}

	f._waitCancelJudge()
	return f.obj
}

func (f *DefaultFuture) Done() <-chan struct{} {
	return f.ctx.Done()
}

func (f *DefaultFuture) Await() Future {
	f.Get()
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
	f.opL.Lock()
	defer f.opL.Unlock()
	f.listeners = append(f.listeners, listener)
	return f
}

func (f *DefaultFuture) self() Future {
	return f
}

func (f *DefaultFuture) Completable() CompletableFuture {
	return f.self().(CompletableFuture)
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
	f.once.Do(func() {
		for _, listener := range f.listeners {
			listener.OperationCompleted(f)
		}
	})
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
	if l.f == nil {
		println("nil f in future listener")
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
	l.Future.(*DefaultFuture).opL.Lock()
	defer l.Future.(*DefaultFuture).opL.Unlock()
	l.Future.(*DefaultFuture).listeners = append(l.Future.(*DefaultFuture).listeners, listener)
	return l
}

func NewFutureListener(f func(f Future)) FutureListener {
	return &_FutureListener{
		Future: NewFuture(nil),
		f:      f,
	}
}
