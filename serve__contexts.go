package syncs

import "context"

type serve_context int

const (
	context_port serve_context = iota
)

func WithServedPort(ctx context.Context, port int) context.Context {
	return context.WithValue(ctx, context_port, port)
}

func ServedPortFrom(ctx context.Context) (int, bool) {
	d := ctx.Value(context_port)
	port, ok := d.(int)
	return port, ok
}
