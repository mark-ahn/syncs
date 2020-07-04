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

type ctx_key_periodic int

const (
	ctx_key_duration ctx_key_periodic = iota
)

func DurationFrom(ctx syncs.Valuable) time.Duration {
	return ctx.Value(ctx_key_duration).(time.Duration)
}

type periodic_print struct{}

func (__ *periodic_print) ServeThread(ctx context.Context, brk syncs.ContextBreakSetter) error {
	th_cnt := syncs.ThreadCounterFrom(ctx)
	ok := th_cnt.AddOrNot(1)
	if !ok {
		return fmt.Errorf("cannot start periodic thread")
	}
	brk.SetValue(ctx_key_duration, time.Second)
	go func() {
		defer th_cnt.Done()
		ticker := time.NewTicker(time.Second)
	loop:
		for {
			select {
			// case <-time.After(time.Second):
			case <-ticker.C:
				fmt.Println(time.Now())
			case <-ctx.Done():
				break loop
			}
		}
	}()

	return nil
}
func Test_Serve(t *testing.T) {
	// ctx, cancel := context.WithCancel(context.TODO())
	// cancel()
	ctx := context.TODO()
	dctx, err := syncs.Serve(ctx, &periodic_print{})
	if err != nil {
		t.Fatal(err)
	}

	du_set := DurationFrom(dctx)
	fmt.Println(du_set)

	ctrl_c := make(chan os.Signal)
	signal.Notify(ctrl_c, os.Interrupt)
	select {
	case <-dctx.Done():
	case <-time.After(3 * du_set):
		dctx.Break(nil)
		<-dctx.Done()
	case <-ctrl_c:
		dctx.Break(nil)
		<-dctx.Done()
	}
	fmt.Println("done", dctx.Err())
}
