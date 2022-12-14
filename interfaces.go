package syncs

import "context"

type Closer interface {
	Close(ctx context.Context) error
}
