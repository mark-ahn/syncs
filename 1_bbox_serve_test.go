package syncs_test

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"testing"
	"time"

	"github.com/mark-ahn/syncs"
)

type periodic_print struct{}

func (__ *periodic_print) ServeThread(ctx context.Context, breaker syncs.ContextBreak) {
	th_cnt := syncs.ThreadCounterFrom(ctx)
	ok := th_cnt.AddOrNot(1)
	if !ok {
		ok = breaker.Break(syncs.ServeFailErrorf("cannot start periodic thread"))
		fmt.Println("bool", ok)
		return
	}
	go func() {
		defer th_cnt.Done()
	loop:
		for {
			select {
			case <-time.After(time.Second):
				fmt.Println(time.Now())
			case <-ctx.Done():
				break loop
			}
		}
	}()
}
func Test_Serve(t *testing.T) {
	// ctx, cancel := context.WithCancel(context.TODO())
	// cancel()
	ctx := context.TODO()
	dctx := syncs.Serve(ctx, &periodic_print{})

	ctrl_c := make(chan os.Signal)
	signal.Notify(ctrl_c, os.Interrupt)
	select {
	case <-dctx.Done():
	case <-ctrl_c:
		dctx.Break(nil)
		<-dctx.Done()
	}
	fmt.Println("done", dctx.Err())
}
