package syncs_test

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/mark-ahn/syncs"
)

func TestThreadCnt(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	go func() {
		<-time.After(time.Microsecond * time.Duration(rand.Int31n(20000)))
		cancel()
	}()
	ctx, done := syncs.WithThreadDoneNotify(ctx, &sync.WaitGroup{})
	cnt := syncs.ThreadCounterFrom(ctx)
	var i int
	for i = 0; i < 1000; i += 1 {
		<-time.After(time.Microsecond)
		cnted := cnt.AddOrNot(1)
		if !cnted {
			break
		}
		go func() {
			<-time.After(time.Second)
			cnt.Done()
		}()
	}
	println(i)
	<-done

}
