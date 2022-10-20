package syncs

import (
	"context"
	"reflect"
)

func ctxkey[T any]() interface{} {
	var zero *T
	return reflect.TypeOf(zero).Elem()
}

func ValueFrom[T any](ctx Valuable) (T, bool) {
	value, ok := ctx.Value(ctxkey[T]()).(T)
	return value, ok
}
func WithValue[T any](ctx context.Context, value T) context.Context {
	return context.WithValue(ctx, ctxkey[T](), value)
}

func SetValue[T any](ctx ValueSetter, value T) (T, bool) {
	d := ctx.SetValue(ctxkey[T](), value)
	old, ok := d.(T)
	return old, ok
}
