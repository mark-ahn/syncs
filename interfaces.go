package syncs

import "context"

type Closer interface {
	Close(ctx context.Context) error
}

var _ Closer = CloseFunc(nil)

type CloseFunc func(ctx context.Context) error

func (__ CloseFunc) Close(ctx context.Context) error {
	return __(ctx)
}
