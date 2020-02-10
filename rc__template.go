package syncs

import "sync/atomic"

type RefCounterOfSome struct {
	cnt    *uint32
	obj    Some
	relase func(Some)
}

func NewRefCounterOfSome(obj Some, release func(Some)) *RefCounterOfSome {
	cnt := uint32(1)
	return &RefCounterOfSome{
		cnt:    &cnt,
		obj:    obj,
		relase: release,
	}
}

func (__ *RefCounterOfSome) Clone() *RefCounterOfSome {
	if atomic.AddUint32(__.cnt, 1) == 1 {
		panic("Release() called during clone")
	}
	return __
}

func (__ *RefCounterOfSome) Object() Some {
	if atomic.LoadUint32(__.cnt) == 0 {
		return nil
	}
	return __.obj
}

func (__ *RefCounterOfSome) Release() {
	if atomic.AddUint32(__.cnt, ^uint32(0)) == 0 {
		__.relase(__.obj)
	}
}
