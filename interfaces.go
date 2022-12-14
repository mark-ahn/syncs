package syncs

import (
	"context"
	"time"
)

type Closer interface {
	Close(ctx context.Context) error
}

var _ Closer = CloseFunc(nil)

type Procedure func(ctx context.Context) error

func (__ Procedure) CallWithTimeout(ctx context.Context, du time.Duration) error {
	return CallWithTimeout(ctx, du, __)
}

type CloseFunc Procedure

func (__ CloseFunc) Close(ctx context.Context) error {
	return __(ctx)
}

func CallWithTimeout(ctx context.Context, du time.Duration, f Procedure) error {
	ctx, stop := context.WithTimeout(ctx, du)
	defer stop()
	return f(ctx)
}
