package syncs

import "sync/atomic"

type RefCounter[Some any] struct {
	cnt     *uint32
	obj     Some
	release func(Some)
}

func NewRefCounter[Some any](obj Some, release func(Some)) *RefCounter[Some] {
	cnt := uint32(1)
	return &RefCounter[Some]{
		cnt:     &cnt,
		obj:     obj,
		release: release,
	}
}

func (__ *RefCounter[Some]) Clone() *RefCounter[Some] {
	if atomic.AddUint32(__.cnt, 1) == 1 {
		panic("Release() called during clone")
	}
	return __
}

func (__ *RefCounter[Some]) Object() (_ Some) {
	if atomic.LoadUint32(__.cnt) == 0 {
		return
	}
	return __.obj
}

func (__ *RefCounter[Some]) Release() {
	if atomic.AddUint32(__.cnt, ^uint32(0)) == 0 {
		__.release(__.obj)
	}
}
