package syncs

import (
	"context"
	"sync"
)

type Valuable interface {
	Value(interface{}) interface{}
}

type ThreadCounter interface {
	Add(int)
	Done()
}

type dummy_thread_counter struct{}

func (_ dummy_thread_counter) Add(int) {}
func (_ dummy_thread_counter) Done()   {}

var dummy_counter = &dummy_thread_counter{}

type context_key int

const (
	context_key_thread_counter context_key = iota
)

func WithThreadCounter(ctx context.Context, counter ThreadCounter) context.Context {
	return context.WithValue(ctx, context_key_thread_counter, counter)
}
func ThreadCounterFrom(ctx Valuable) ThreadCounter {
	v, ok := ctx.Value(context_key_thread_counter).(ThreadCounter)
	if !ok {
		return dummy_counter
	}
	return v
}

type WaitGroup interface {
	ThreadCounter
	Wait()
}

type cnt_starter struct {
	group   WaitGroup
	starter func()
}

func new_cnt_starter(group WaitGroup, f func()) *cnt_starter {
	once := sync.Once{}
	return &cnt_starter{
		group: group,
		starter: func() {
			once.Do(f)
		},
	}
}
func (__ *cnt_starter) Add(i int) {
	__.starter()
	__.group.Add(i)
}

func (__ *cnt_starter) Done() { __.group.Done() }

// func (__ *cnt_starter) Wait() { __.group.Wait() }

func WithThreadDoneNotify(ctx context.Context, threads *sync.WaitGroup) (context.Context, <-chan struct{}) {
	p_cnt := ThreadCounterFrom(ctx)

	start_ctx, cnt_start := context.WithCancel(ctx)

	in_ctx := WithThreadCounter(ctx, new_cnt_starter(threads, cnt_start))
	done_ch := make(chan struct{})
	p_cnt.Add(1)
	go func() {
		defer p_cnt.Done()
		defer close(done_ch)
		defer threads.Wait()

	loop:
		for {
			select {
			case <-start_ctx.Done():
				// consider Add() is not called, but start_ctx.Done() because of parent context
				// cnt_start should be called at least onece to avoid resource leak
				cnt_start()
				break loop
			}
		}
	}()
	return in_ctx, done_ch
}
