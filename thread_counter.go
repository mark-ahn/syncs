package syncs

import (
	"context"
	"sync"
)

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
func ThreadCounterFrom(ctx context.Context) ThreadCounter {
	v, ok := ctx.Value(context_key_thread_counter).(ThreadCounter)
	if !ok {
		return dummy_counter
	}
	return v
}

func WithThreadDoneNotify(ctx context.Context, threads *sync.WaitGroup) (context.Context, <-chan struct{}) {
	ctx = WithThreadCounter(ctx, threads)
	done_ch := make(chan struct{})
	go func() {
		defer close(done_ch)
		defer func() {
			threads.Wait()
			// if release != nil {
			// 	release()
			// }
		}()

	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			}
		}
	}()
	return ctx, done_ch
}
