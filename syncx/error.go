package syncx

import "sync/atomic"

type AtomicError struct {
	err atomic.Value
}

func (ae *AtomicError) Set(err error) {
	if err != nil {
		ae.err.Store(err)
	}
}

func (ae *AtomicError) Load() error {
	if v := ae.err.Load(); v != nil {
		return v.(error)
	}

	return nil
}
