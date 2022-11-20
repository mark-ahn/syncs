package syncs

import (
	"context"
	"sync"

	"github.com/mark-ahn/metrics"
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
	TryAdd(int) error
	Done()
}

func try_add_thread_count(tcnt ThreadCounter, n int) error {
	ok := tcnt.AddOrNot(n)
	if !ok {
		return TrySpawnThreadOnContextDoneErrorf("with count %d", n)
	}
	return nil
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
func (__ dummy_counter) TryAdd(n int) error {
	return try_add_thread_count(__, n)
}
func (dummy_counter) Done() {}

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

// --------------------------------------------------------------------------------------

type WaitGroup interface {
	SyncCounter
	Wait()
}

type cnt_starter struct {
	counter SyncCounter
	done    <-chan struct{}
	mutext  *sync.Mutex
	scope   Scope
}

func new_cnt_starter(group SyncCounter, mutext *sync.Mutex, done <-chan struct{}, scope Scope) *cnt_starter {
	return &cnt_starter{
		counter: group,
		done:    done,
		mutext:  mutext,
		scope:   scope,
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
		__.scope.PutMetric(ThreadCountMetric{Delta: i}, nil)
		return true
	}
}

func (__ *cnt_starter) Done() {
	__.counter.Done()
	__.scope.PutMetric(ThreadCountMetric{Delta: -1}, nil)
}

func (__ *cnt_starter) TryAdd(n int) error {
	return try_add_thread_count(__, n)
}

// --------------------------------------------------------------------------------------

// WithThreadDoneNotify inserts thread-counter into context, then returns a channel which would be closed
// after all threads spawnded with the thread-counter are terminated.
func WithThreadDoneNotify(ctx context.Context, threads WaitGroup) (context.Context, <-chan struct{}) {
	p_cnt := ThreadCounterFrom(ctx)

	// sync_ch & mutex confirms that threads.Add() always be called before thread.Wait()
	// in other words it makes avoid threads.Add() be called after thread.Wait()
	mutex := &sync.Mutex{}
	sync_ch := make(chan struct{})

	scope := metrics.ScopeFromOrDummy[MetricData](ctx)

	in_ctx := WithThreadCounter(ctx, new_cnt_starter(threads, mutex, sync_ch, scope))
	// done_ch is closed after all threads are terminated
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

		<-ctx.Done()
		mutex.Lock()
		close(sync_ch)
		mutex.Unlock()
	}()
	return in_ctx, done_ch
}
