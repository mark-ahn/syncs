package syncs

import (
	"context"
	"sync"
)

type ThreadServer interface {
	ServeThread(ctx context.Context, breaker ContextBreak)
}

func Serve(ctx context.Context, server ThreadServer) SrvHandle {
	in_ctx, cancel := context.WithCancel(ctx)
	in_ctx, done := WithThreadDoneNotify(in_ctx, &sync.WaitGroup{})
	rctx := NewDoneChContext(in_ctx, done, cancel)
	go func() {
		defer cancel()
		<-rctx.Done()
	}()

	server.ServeThread(in_ctx, rctx)

	return rctx
}
