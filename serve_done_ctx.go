package syncs

import (
	"context"
	"sync"
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

type ValueSetter interface {
	SetValue(interface{}, interface{}) interface{}
}

type ContextBreakSetter interface {
	ContextBreak
	ValueSetter
}

type ServeHandle interface {
	contextDoneErr
	ContextBreak
	Valuable
}

type DoneChContext struct {
	pctx   context.Context
	done   <-chan struct{}
	err    error
	set    int32
	cancel func()

	values     map[interface{}]interface{}
	value_lock sync.RWMutex
}

func NewDoneChContext(pctx context.Context, done <-chan struct{}, cancel func()) *DoneChContext {
	__ := &DoneChContext{
		pctx:   pctx,
		done:   done,
		err:    nil,
		set:    0,
		cancel: cancel,

		values:     make(map[interface{}]interface{}),
		value_lock: sync.RWMutex{},
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

func (__ *DoneChContext) Value(k interface{}) interface{} {
	__.value_lock.RLock()
	defer __.value_lock.RUnlock()
	return __.values[k]
}
func (__ *DoneChContext) SetValue(k, v interface{}) interface{} {
	__.value_lock.Lock()
	old := __.values[k]
	__.values[k] = v
	__.value_lock.Unlock()
	return old
}
