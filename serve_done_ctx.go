package syncs

import (
	"context"
	"sync/atomic"
)

type contextDoneErr interface {
	contextDone
	contextErr
}

type contextDone interface {
	Done() <-chan struct{}
}

type contextErr interface {
	Err() error
}

type ContextBreak interface {
	Break(error) bool
}

type SrvHandle interface {
	contextDoneErr
	ContextBreak
}

type DoneChContext struct {
	pctx   context.Context
	done   <-chan struct{}
	err    error
	set    int32
	cancel func()
}

func NewDoneChContext(pctx context.Context, done <-chan struct{}, cancel func()) *DoneChContext {
	__ := &DoneChContext{
		pctx:   pctx,
		done:   done,
		err:    nil,
		set:    0,
		cancel: cancel,
	}
	return __
}

func (__ *DoneChContext) Break(err error) bool {
	select {
	case <-__.done:
		return false
	default:
		if atomic.AddInt32(&__.set, 1) != 1 {
			atomic.AddInt32(&__.set, -1)
			return false
		}
		__.err = err
		__.cancel()
		return true
	}
}

func (__ *DoneChContext) Done() <-chan struct{} {
	return __.done
}

func (__ *DoneChContext) Err() error {
	select {
	case <-__.done:
		switch {
		case atomic.LoadInt32(&__.set) == 0:
			return __.pctx.Err()
		default:
			return __.err
		}
	default:
		return nil
	}
}
