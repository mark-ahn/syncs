package syncs_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/mark-ahn/syncs"
)

func TestThreadCnt(t *testing.T) {
	ctx, done := syncs.WithThreadDoneNotify(context.TODO(), &sync.WaitGroup{})
	cnt := syncs.ThreadCounterFrom(ctx)
	cnt.Add(1)
	go func() {
		<-time.After(time.Second)
		cnt.Done()
	}()
	<-done

}
