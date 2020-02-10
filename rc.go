package syncs

// type RefCounter interface {
// 	Clone() RefCounter
// 	Interface() interface{}
// 	Release()
// }

// type ref_counter struct {
// 	cnt    *uint32
// 	obj    interface{}
// 	relase func(interface{})
// }

// func NewRefCounter(obj interface{}, release func(interface{})) RefCounter {
// 	cnt := uint32(1)
// 	return &ref_counter{
// 		cnt:    &cnt,
// 		obj:    obj,
// 		relase: release,
// 	}
// }

// func (__ *ref_counter) Clone() RefCounter {
// 	if atomic.AddUint32(__.cnt, 1) == 1 {
// 		panic("Release() called during clone")
// 	}
// 	return __
// }

// func (__ *ref_counter) Interface() interface{} {
// 	if atomic.LoadUint32(__.cnt) == 0 {
// 		return nil
// 	}
// 	return __.obj
// }

// func (__ *ref_counter) Release() {
// 	if atomic.AddUint32(__.cnt, ^uint32(0)) == 0 {
// 		__.relase(__.obj)
// 	}
// }
