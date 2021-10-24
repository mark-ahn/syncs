package syncs

type serve_context int

const (
	context_port serve_context = iota
)

func SetServedPort(ctx valueSetter, port int) (int, bool) {
	d := ctx.SetValue(context_port, port)
	port, ok := d.(int)
	return port, ok
}

func ServedPortFrom(ctx Valuable) (int, bool) {
	d := ctx.Value(context_port)
	port, ok := d.(int)
	return port, ok
}
