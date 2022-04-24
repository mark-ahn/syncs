package syncs

import "sync"

type TypedPool[T any] struct {
	// alloc func() T
	free func(T)
	pool sync.Pool
}

func NewTypedPool[T any](alloc func() T, free func(T)) *TypedPool[T] {
	return &TypedPool[T]{
		// alloc: alloc,
		free: free,
		pool: sync.Pool{
			New: func() any { return alloc() },
		},
	}
}

func (__ *TypedPool[T]) Get() T {
	return __.pool.Get().(T)
}

func (__ *TypedPool[T]) Put(d T) {
	__.free(d)
	__.pool.Put(d)
}
