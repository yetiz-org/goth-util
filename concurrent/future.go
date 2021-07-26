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
	f.init()
	return f
}

func NewSucceededFuture(obj interface{}) Future {
	f := &DefaultFuture{}
	f.ctx, f.cancel = context.WithCancel(context.Background())
	f.init()
	f.Complete(obj)
	return f
}

func NewCancelledFuture() Future {
	f := &DefaultFuture{}
	f.ctx, f.cancel = context.WithCancel(context.Background())
	f.init()
	f.Cancel()
	return f
}

func NewFailedFuture(err error) Future {
	f := &DefaultFuture{}
	f.ctx, f.cancel = context.WithCancel(context.Background())
	f.init()
	f.Fail(err)
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

func (f *DefaultFuture) init() {
	f.opL.Lock()
	go func() {
		f.opL.Unlock()
		<-f.ctx.Done()
		f._cancelJudge()
	}()

	f.opL.Lock()
	f.opL.Unlock()
}

func (f *DefaultFuture) _cancelJudge() {
	if atomic.CompareAndSwapInt32(&f.state, stateWait, stateCancel) {
		f.err = f.ctx.Err()
		f.callListener()
	}
}

func (f *DefaultFuture) Get() interface{} {
	if f.IsDone() {
		return f.obj
	}

	<-f.ctx.Done()
	f._cancelJudge()
	return f.obj
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
	Future
	f func(f Future)
}

func (l *_FutureListener) OperationCompleted(f Future) {
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
