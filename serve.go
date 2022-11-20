package syncs

import (
	"context"
	"sync"

	"github.com/mark-ahn/goes"
	"github.com/mark-ahn/metrics"
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

type ServeOpts struct {
	Id string
}

func DefaultServeOpts() *ServeOpts {
	return &ServeOpts{}
}

type ServeOpt = goes.OptionSetter[ServeOpts]

// Serve serves all threads from server, then returns ServeContext which is a handle
// that can signal all threads to be terminated using Break() & confirms all threads are terminated using Done()
func Serve(ctx context.Context, server ThreadServer, opts ...ServeOpt) (ServeContext, error) {
	opt := goes.ApplyOptions(DefaultServeOpts(), opts...)

	in_ctx, stop := context.WithCancel(ctx)

	if opt.Id != "" {
		in_ctx, _ = metrics.OverrideScopeWithLabelOr[Scope](ctx, []string{opt.Id})
	}
	in_ctx, done := WithThreadDoneNotify(in_ctx, &sync.WaitGroup{})
	rctx := NewDoneChContext(in_ctx, done, stop)
	go func() {
		defer stop()
		<-rctx.Done()
	}()

	err := server.ServeThread(in_ctx, rctx)
	if err != nil {
		stop()
		return nil, err
	}

	return rctx, nil
}
