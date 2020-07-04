package syncs

import (
	"context"
	"sync"
)

type ThreadServer interface {
	ServeThread(ctx context.Context, breaker ContextBreakSetter) error
}

type ThreadServerFunc func(context.Context, ContextBreakSetter)

func (__ ThreadServerFunc) ServeThread(ctx context.Context, bk ContextBreakSetter) { __(ctx, bk) }

func Serve(ctx context.Context, server ThreadServer) (ServeHandle, error) {
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
