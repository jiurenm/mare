package singleflight

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
)

// errGoexit indicates the runtime.Goexit was called in
// the user given function.
var errGoexit = errors.New("runtime.Goexit was called")

type panicError struct {
	value any
	stack []byte
}

// Error implements error interface.
func (p *panicError) Error() string {
	return fmt.Sprintf("%v\n\n%s", p.value, p.stack)
}

func newPanicError(v any) error {
	stack := debug.Stack()

	// The first line of the stack trace is of the form "goroutine N [status]:"
	// but by the time the panic reaches Do the goroutine may no longer exist
	// and its status will have changed. Trim out the misleading line.
	if line := bytes.IndexByte(stack[:], '\n'); line >= 0 {
		stack = stack[line+1:]
	}
	return &panicError{value: v, stack: stack}
}

type call[T any] struct {
	val T
	err error

	chans []chan<- Result[T]
	wg    sync.WaitGroup

	dups int

	forgotten bool
}

type Group[T any] struct {
	m  map[string]*call[T]
	mu sync.Mutex
}

type Result[T any] struct {
	Val    T
	Err    error
	Shared bool
}

func (g *Group[T]) Do(key string, fn func() (T, error)) (v T, err error, shared bool) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call[T])
	}
	if c, ok := g.m[key]; ok {
		c.dups++
		g.mu.Unlock()
		c.wg.Wait()

		if e, ok := c.err.(*panicError); ok {
			panic(e)
		} else if c.err == errGoexit {
			runtime.Goexit()
		}
		return c.val, c.err, true
	}
	c := new(call[T])
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	g.doCall(c, key, fn)
	return c.val, c.err, c.dups > 0
}

func (g *Group[T]) DoChan(key string, fn func() (T, error)) <-chan Result[T] {
	ch := make(chan Result[T], 1)
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call[T])
	}
	if c, ok := g.m[key]; ok {
		c.dups++
		c.chans = append(c.chans, ch)
		g.mu.Unlock()
		return ch
	}
	c := &call[T]{chans: []chan<- Result[T]{ch}}
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	go g.doCall(c, key, fn)
	return ch
}

func (g *Group[T]) doCall(c *call[T], key string, fn func() (T, error)) {
	normalReturn := false
	recoverd := false

	defer func() {
		if !normalReturn && !recoverd {
			c.err = errGoexit
		}

		c.wg.Done()
		g.mu.Lock()
		defer g.mu.Unlock()
		if !c.forgotten {
			delete(g.m, key)
		}

		if e, ok := c.err.(*panicError); ok {
			if len(c.chans) > 0 {
				go panic(e)
				select {}
			} else {
				panic(e)
			}
		} else if c.err == errGoexit {

		} else {
			for _, ch := range c.chans {
				ch <- Result[T]{c.val, c.err, c.dups > 0}
			}
		}

	}()

	func() {
		defer func() {
			if !normalReturn {
				if r := recover(); r != nil {
					c.err = newPanicError(r)
				}
			}
		}()

		c.val, c.err = fn()
		normalReturn = true
	}()

	if !normalReturn {
		recoverd = true
	}
}

func (g *Group[T]) Forget(key string) {
	g.mu.Lock()
	if c, ok := g.m[key]; ok {
		c.forgotten = true
	}
	delete(g.m, key)
	g.mu.Unlock()
}
