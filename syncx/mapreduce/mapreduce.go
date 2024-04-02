package mr

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/jiurenm/mare/syncx"
)

const (
	defaultWorkers = 16
	minWorkers     = 1
)

var (
	// ErrCancelWithNil is an error that mapreduce was cancelled with nil.
	ErrCancelWithNil = errors.New("mapreduce cancelled with nil")
	// ErrReduceNoOutput is an error that reduce did not output a value.
	ErrReduceNoOutput = errors.New("reduce not writing value")
)

type (
	generateFunc[T any]       func(source chan<- T)
	mapFunc[T any, K any]     func(item T, writer Writer[K])
	mapperFunc[T any, K any]  func(item T, writer Writer[K], cancel func(error))
	reducerFunc[T any, K any] func(pipe <-chan T, writer Writer[K], cancel func(error))
	voidReducerFunc           func(pipe <-chan any, cancel func(error))

	Option func(opts *mapReduceOptions)

	mapperContext[T, K any] struct {
		ctx       context.Context
		mapper    func(item T, writer Writer[K])
		source    <-chan T
		panicChan *onceChan
		collector chan<- K
		doneChan  <-chan struct{}
		workers   int
	}

	mapReduceOptions struct {
		ctx     context.Context
		workers int
	}

	Writer[T any] interface {
		Write(v T)
	}
)

type onceChan struct {
	channel chan any
	wrote   int32
}

type DoneChan struct {
	done chan struct{}
	once sync.Once
}

func NewDoneChan() *DoneChan {
	return &DoneChan{
		done: make(chan struct{}),
	}
}

func (dc *DoneChan) Close() {
	dc.once.Do(func() {
		close(dc.done)
	})
}

func (dc *DoneChan) Done() chan struct{} {
	return dc.done
}

func (oc *onceChan) write(val any) {
	if atomic.CompareAndSwapInt32(&oc.wrote, 0, 1) {
		oc.channel <- val
	}
}

func Finish(fns ...func() error) error {
	if len(fns) == 0 {
		return nil
	}

	return mapReduceVoid(func(source chan<- any) {
		for _, fn := range fns {
			source <- fn
		}
	}, func(item any, writer Writer[any], cancel func(error)) {
		fn := item.(func() error)
		if err := fn(); err != nil {
			cancel(err)
		}
	}, func(pipe <-chan any, cancel func(error)) {
		drain(pipe)
	}, WithWorkers(len(fns)))
}

func mapReduceVoid(generate func(source chan<- any), mapper func(item any, writer Writer[any], cancel func(error)), reducer voidReducerFunc, opts ...Option) error {
	_, err := MapReduce(generate, mapper, func(input <-chan any, writer Writer[any], cancel func(error)) {
		reducer(input, cancel)
		writer.Write(struct{}{})
	}, opts...)

	return err
}

func WithWorkers(workers int) Option {
	return func(opts *mapReduceOptions) {
		if workers < minWorkers {
			opts.workers = minWorkers
		} else {
			opts.workers = workers
		}
	}
}

func MapReduce[T any, K any, R any](generate generateFunc[T], mapper mapperFunc[T, K], reducer reducerFunc[K, R], opts ...Option) (R, error) {

	panicChan := &onceChan{channel: make(chan any)}
	source := buildSource(generate, panicChan)

	return MapReduceWithPanicChan(source, panicChan, mapper, reducer, opts...)
}

func MapReduceWithSource[T any, K any, R any](source <-chan T, mapper mapperFunc[T, K], reducer reducerFunc[K, R], opts ...Option) (R, error) {
	options := buildOptions(opts...)
	output := make(chan R)

	defer func() {
		for range output {
			panic("more than one element written in reducer")
		}
	}()

	collector := make(chan K, options.workers)
	done := NewDoneChan()
	writer := newGuardedWriter(context.Background(), output, done.Done())

	var (
		closeOnce sync.Once
		retErr    syncx.AtomicError
	)

	finish := func() {
		closeOnce.Do(func() {
			done.Close()
			close(output)
		})
	}
	cancel := once(func(err error) {
		if err != nil {
			retErr.Set(err)
		} else {
			retErr.Set(ErrCancelWithNil)
		}

		drain(source)
		finish()
	})

	go func() {
		defer func() {
			drain(collector)

			if r := recover(); r != nil {
				cancel(fmt.Errorf("%v", r))
			} else {
				finish()
			}
		}()

		reducer(collector, writer, cancel)
	}()

	go executeMappers(func(item T, w Writer[K]) {
		mapper(item, w, cancel)
	}, source, collector, done.Done(), options.workers)

	value, ok := <-output

	if err := retErr.Load(); err != nil {
		var ret R

		return ret, err
	} else if ok {
		return value, nil
	} else {
		var ret R

		return ret, ErrReduceNoOutput
	}
}

func MapReduceWithPanicChan[T, K, R any](source <-chan T, panicChan *onceChan, mapper mapperFunc[T, K], reducer reducerFunc[K, R], opts ...Option) (val R, err error) {
	options := buildOptions(opts...)

	output := make(chan R)
	defer func() {
		for range output {
			panic("more than one element written in reducer")
		}
	}()

	collector := make(chan K, options.workers)
	done := make(chan struct{})
	writer := newGuardedWriter(options.ctx, output, done)
	var closeOnce sync.Once

	var retErr syncx.AtomicError
	finish := func() {
		closeOnce.Do(func() {
			close(done)
			close(output)
		})
	}
	cancel := once(func(err error) {
		if err != nil {
			retErr.Set(err)
		} else {
			retErr.Set(ErrCancelWithNil)
		}

		drain(source)
		finish()
	})

	go func() {
		defer func() {
			drain(collector)
			if r := recover(); r != nil {
				panicChan.write(r)
			}
			finish()
		}()

		reducer(collector, writer, cancel)
	}()

	go executeMappersWithCtx(mapperContext[T, K]{
		ctx: options.ctx,
		mapper: func(item T, w Writer[K]) {
			mapper(item, w, cancel)
		},
		source:    source,
		panicChan: panicChan,
		collector: collector,
		doneChan:  done,
		workers:   options.workers,
	})

	select {
	case <-options.ctx.Done():
		cancel(context.DeadlineExceeded)
		err = context.DeadlineExceeded
	case v := <-panicChan.channel:
		drain(source)
		panic(v)
	case v, ok := <-output:
		if e := retErr.Load(); e != nil {
			err = e
		} else if ok {
			val = v
		} else {
			err = ErrReduceNoOutput
		}
	}

	return
}

func executeMappers[T any, K any](mapper mapFunc[T, K], input <-chan T, collector chan<- K, done <-chan struct{}, workers int) {
	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
		close(collector)
	}()

	pool := make(chan struct{}, workers)
	writer := newGuardedWriter(context.Background(), collector, done)

	for {
		select {
		case <-done:
			return
		case pool <- struct{}{}:
			item, ok := <-input
			if !ok {
				<-pool

				return
			}

			wg.Add(1)
			syncx.GoSafe(func() {
				defer func() {
					wg.Done()
					<-pool
				}()
				mapper(item, writer)
			})
		}
	}
}

func executeMappersWithCtx[T, K any](mCtx mapperContext[T, K]) {
	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
		close(mCtx.collector)
		drain(mCtx.source)
	}()

	var failed int32
	pool := make(chan struct{}, mCtx.workers)
	writer := newGuardedWriter(mCtx.ctx, mCtx.collector, mCtx.doneChan)
	for atomic.LoadInt32(&failed) == 0 {
		select {
		case <-mCtx.ctx.Done():
			return
		case <-mCtx.doneChan:
			return
		case pool <- struct{}{}:
			item, ok := <-mCtx.source
			if !ok {
				<-pool
				return
			}

			wg.Add(1)
			go func() {
				defer func() {
					if r := recover(); r != nil {
						atomic.AddInt32(&failed, 1)
						mCtx.panicChan.write(r)
					}
					wg.Done()
					<-pool
				}()

				mCtx.mapper(item, writer)
			}()
		}
	}
}

func buildSource[T any](generate generateFunc[T], panicChan *onceChan) chan T {
	source := make(chan T)

	syncx.GoSafe(func() {
		defer func() {
			if r := recover(); r != nil {
				panicChan.write(r)
			}
			close(source)
		}()
		generate(source)
	})

	return source
}

func buildOptions(opts ...Option) *mapReduceOptions {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}

	return options
}

func newOptions() *mapReduceOptions {
	return &mapReduceOptions{workers: defaultWorkers}
}

func once(fn func(error)) func(error) {
	once := new(sync.Once)

	return func(err error) {
		once.Do(func() {
			fn(err)
		})
	}
}

func drain[T any](channel <-chan T) {
	for range channel {
	}
}

type guardedWriter[T any] struct {
	ctx     context.Context
	channel chan<- T
	done    <-chan struct{}
}

func newGuardedWriter[T any](ctx context.Context, channel chan<- T, done <-chan struct{}) guardedWriter[T] {
	return guardedWriter[T]{
		ctx:     ctx,
		channel: channel,
		done:    done,
	}
}

func (gw guardedWriter[T]) Write(v T) {
	select {
	case <-gw.done:
		return
	default:
		gw.channel <- v
	}
}
