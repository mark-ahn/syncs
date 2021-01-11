package syncs

import (
	"context"
	"log"
)

func ThreadCounterWith(ctx context.Context, cnt int) (ThreadCounter, error) {
	th_cnt := ThreadCounterFrom(ctx)
	ok := th_cnt.AddOrNot(1)
	if !ok {
		return nil, ServeThreadErrorf("context done")
	}
	return th_cnt, nil
}

func NewThreadServerFunc(f func(done <-chan struct{}) error, breakOnTerminated bool) ThreadServerFunc {
	return func(ctx context.Context, tctx ThreadContext) error {
		th_cnt, err := ThreadCounterWith(ctx, 1)
		if err != nil {
			return err
		}
		go func() {
			defer th_cnt.Done()
			err := f(ctx.Done())
			switch {
			case breakOnTerminated:
				tctx.Break(err)
			}
		}()
		return nil
	}
}

func NewThreadServerFuncLoop(ch <-chan interface{}, f func(d interface{}, ok bool) error, breakOnTerminated bool) ThreadServerFunc {
	return NewThreadServerFunc(func(done <-chan struct{}) error {
	loop:
		for {
			select {
			case <-done:
				break loop
			case d, ok := <-ch:
				err := f(d, ok)
				if err != nil {
					log.Printf("terminating thread: %v", err)
					break loop
				}
			}
		}
		return nil
	}, breakOnTerminated)
}
