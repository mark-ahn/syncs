// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package syncs

import "sync/atomic"

type RefCounterOfInterface struct {
	cnt    *uint32
	obj    interface{}
	relase func(interface{})
}

func NewRefCounterOfInterface(obj interface{}, release func(interface{})) *RefCounterOfInterface {
	cnt := uint32(1)
	return &RefCounterOfInterface{
		cnt:    &cnt,
		obj:    obj,
		relase: release,
	}
}

func (__ *RefCounterOfInterface) Clone() *RefCounterOfInterface {
	if atomic.AddUint32(__.cnt, 1) == 1 {
		panic("Release() called during clone")
	}
	return __
}

func (__ *RefCounterOfInterface) Object() interface{} {
	if atomic.LoadUint32(__.cnt) == 0 {
		return nil
	}
	return __.obj
}

func (__ *RefCounterOfInterface) Release() {
	if atomic.AddUint32(__.cnt, ^uint32(0)) == 0 {
		__.relase(__.obj)
	}
}
