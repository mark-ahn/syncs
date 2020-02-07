package syncs_test

import (
	"fmt"
	"io"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/mark-ahn/syncs"
)

type dummy_stream struct {
	readCnt  int32
	writeCnt int32
}

func (__ *dummy_stream) Read(p []byte) (n int, err error) {
	fmt.Printf("read\n")
	atomic.AddInt32(&__.readCnt, 1)
	return 0, nil
}

func (__ *dummy_stream) Write(p []byte) (n int, err error) {
	fmt.Printf("write\n")
	atomic.AddInt32(&__.writeCnt, 1)
	return 0, nil
}

func (__ *dummy_stream) Close() error {
	fmt.Printf("closed\n")
	return nil
}

func init() {
	rand.Seed(time.Now().Unix())
}
func TestRc(t *testing.T) {
	rc := syncs.NewRefCounter(&dummy_stream{}, func(obj interface{}) {
		obj.(io.Closer).Close()
	})

	defer rc.Release()

	threads := &sync.WaitGroup{}

	var some_work func(syncs.RefCounter)
	some_work = func(rc syncs.RefCounter) {
		defer func() {
			rc.Release()
			threads.Done()
		}()

		if rand.Intn(10) == 0 {
			threads.Add(1)
			go some_work(rc.Clone())
		}

		stream := rc.Interface().(io.ReadWriter)

		stream.Read([]byte{})
		stream.Write([]byte{})

	}

	for i := 0; i < 100; i += 1 {
		threads.Add(1)
		go some_work(rc.Clone())
	}

	threads.Wait()
	ds := rc.Interface().(*dummy_stream)
	fmt.Printf("done read: %v, write %v\n", ds.readCnt, ds.writeCnt)
}
