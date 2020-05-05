package syncs

import (
	"context"
	"sync"
)

type Valuable interface {
	Value(interface{}) interface{}
}

type SyncCounter interface {
	Add(int)
	Done()
}

type ThreadCounter interface {
	AddOrNot(int) bool
	Done()
}

type dummy_counter <-chan struct{}

func (__ dummy_counter) AddOrNot(int) bool {
	select {
	case <-__:
		return false
	default:
		return true
	}
}
func (_ dummy_counter) Done() {}

type context_key int

const (
	context_key_thread_counter context_key = iota
)

func WithThreadCounter(ctx context.Context, counter ThreadCounter) context.Context {
	return context.WithValue(ctx, context_key_thread_counter, counter)
}

func ThreadCounterFrom(ctx context.Context) ThreadCounter {
	v, ok := ctx.Value(context_key_thread_counter).(ThreadCounter)
	if !ok {
		return dummy_counter(ctx.Done())
	}
	return v
}

type WaitGroup interface {
	SyncCounter
	Wait()
}

type cnt_starter struct {
	counter SyncCounter
	done    <-chan struct{}
	mutext  *sync.Mutex
}

func new_cnt_starter(group SyncCounter, mutext *sync.Mutex, done <-chan struct{}) *cnt_starter {
	return &cnt_starter{
		counter: group,
		done:    done,
		mutext:  mutext,
	}
}
func (__ *cnt_starter) AddOrNot(i int) bool {
	__.mutext.Lock()
	defer __.mutext.Unlock()

	select {
	case <-__.done:
		return false
	default:
		__.counter.Add(i)
		return true
	}
}

func (__ *cnt_starter) Done() {
	__.counter.Done()
}

// func (__ *cnt_starter) Wait() { __.group.Wait() }

func WithThreadDoneNotify(ctx context.Context, threads WaitGroup) (context.Context, <-chan struct{}) {
	p_cnt := ThreadCounterFrom(ctx)

	mutex := &sync.Mutex{}
	sync_ch := make(chan struct{})

	in_ctx := WithThreadCounter(ctx, new_cnt_starter(threads, mutex, sync_ch))
	done_ch := make(chan struct{})

	cnted := p_cnt.AddOrNot(1)
	if !cnted {
		close(done_ch)
		close(sync_ch)
		return in_ctx, done_ch
	}

	go func() {
		defer p_cnt.Done()
		defer close(done_ch)
		defer threads.Wait()

	loop:
		for {
			select {
			// case <-start_ctx.Done():
			// 	break loop
			case <-ctx.Done():
				mutex.Lock()
				close(sync_ch)
				mutex.Unlock()
				break loop
			}
		}
	}()
	return in_ctx, done_ch
}
