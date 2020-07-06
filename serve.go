package syncs

import (
	"context"
	"sync"
)

type ThreadServer interface {
	ServeThread(ctx context.Context, tctx ThreadContext) error
}

type ThreadServerFunc func(context.Context, ThreadContext) error

func (__ ThreadServerFunc) ServeThread(ctx context.Context, tctx ThreadContext) error {
	return __(ctx, tctx)
}

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
		return nil, err
	}

	return rctx, nil
}
