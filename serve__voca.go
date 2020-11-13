package syncs

import (
	"context"
	"log"
)

// func (__ *) ServeThread(ctx context.Context, tctx syncs.ThreadContext) error {
//   th_cnt := syncs.ThreadCounterFrom(ctx)
//   ok := th_cnt.AddOrNot(1)
//   if !ok {
// 	return fmt.Errorf("cannot serve thread: context done")
//   }

//   go func() {
// 	defer th_cnt.Done()
//   loop:
// 	for {
// 	  select {
// 	  case <-ctx.Done():
// 		break loop
// 	  }
// 	}
//   }()
//   return nil
// }

func NewThreadServerFuncLoop(ch <-chan interface{}, f func(d interface{}, ok bool) error) ThreadServerFunc {
	return func(ctx context.Context, tctx ThreadContext) error {
		th_cnt := ThreadCounterFrom(ctx)
		ok := th_cnt.AddOrNot(1)
		if !ok {
			return ServeThreadErrorf("context done")
		}

		go func() {
			defer th_cnt.Done()
		loop:
			for {
				select {
				case <-ctx.Done():
					break loop
				case d, ok := <-ch:
					err := f(d, ok)
					if err != nil {
						log.Printf("terminating thread: %v", err)
						break loop
					}
				}
			}
		}()
		return nil
	}
}
