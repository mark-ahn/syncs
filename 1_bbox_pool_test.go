package syncs_test

import (
	"testing"

	"github.com/mark-ahn/syncs"
)

func Test_TypedPool(t *testing.T) {
	type some_composte struct {
		uptr *uint
	}
	pool := syncs.NewTypedPool(func() *some_composte {
		var d uint
		return &some_composte{uptr: &d}
	}, func(d *some_composte) {
		d.uptr = nil
	})
	storage := make([]*some_composte, 1000)
	for i := 0; i < len(storage); i += 1 {
		storage[i] = pool.Get()
	}
	for _, d := range storage {
		pool.Put(d)
	}
}
