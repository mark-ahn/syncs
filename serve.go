package syncs

import (
	"context"
	"sync"
)

type ThreadServer interface {
	// ctx provides normal context
	// tctx provides breakable handle which signals all threads that are spawned together to be terminated,
	// also provides shared storage cross threads using Value/SetValue interface.
	ServeThread(ctx context.Context, tctx ThreadContext) error
}

type ThreadServerFunc func(context.Context, ThreadContext) error

func (__ ThreadServerFunc) ServeThread(ctx context.Context, tctx ThreadContext) error {
	return __(ctx, tctx)
}

// Serve serves all threads from server, then returns ServeContext which is a handle
// that can signal all threads to be terminated using Break() & confirms all threads are terminated using Done()
func Serve(ctx context.Context, server ThreadServer) (ServeContext, error) {
	in_ctx, cancel := context.WithCancel(ctx)
	in_ctx, done := WithThreadDoneNotify(in_ctx, &sync.WaitGroup{})
	rctx := NewDoneChContext(in_ctx, done, cancel)
	go func() {
		defer cancel()
		<-rctx.Done()
	}()

	err := server.ServeThread(in_ctx, rctx)
	if err != nil {
		cancel()
		return nil, err
	}

	return rctx, nil
}
