package syncs

import "sync/atomic"

type RefCounter[T any] struct {
	cnt    *uint32
	obj    T
	relase func(T)
}

func NewRefCounter[T any](obj T, release func(T)) *RefCounter[T] {
	cnt := uint32(1)
	return &RefCounter[T]{
		cnt:    &cnt,
		obj:    obj,
		relase: release,
	}
}

func (__ *RefCounter[T]) Clone() *RefCounter[T] {
	if atomic.AddUint32(__.cnt, 1) == 1 {
		panic("Release() called during clone")
	}
	return __
}

func (__ *RefCounter[T]) Item() T {
	if atomic.LoadUint32(__.cnt) == 0 {
		var empty T
		return empty
	}
	return __.obj
}

func (__ *RefCounter[T]) Release() {
	if atomic.AddUint32(__.cnt, ^uint32(0)) == 0 {
		__.relase(__.obj)
	}
}
