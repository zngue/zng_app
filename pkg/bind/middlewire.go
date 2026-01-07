package bind

import "context"

type MiddlewareFn = func(ctx context.Context) (tx context.Context, err error)

var middlewares []MiddlewareFn

func SetMiddleWires(fns ...MiddlewareFn) {
	middlewares = fns
}
func GetMiddleWires(ctx context.Context) (bx context.Context, err error) {
	if len(middlewares) > 0 {
		for _, fn := range middlewares {
			var tx context.Context
			tx, err = fn(ctx)
			if err != nil {
				return
			}
			if tx != nil {
				ctx = tx
			}
		}
	}
	bx = ctx
	return
}
func ClearMiddleWires() {
	middlewares = nil
}

// Handler defines the handler invoked by Middleware.
type Handler func(ctx context.Context, req any) (rs any, err error)

// Middleware is HTTP/gRPC transport middleware.
type Middleware func(Handler) Handler

// Chain returns a Middleware that specifies the chained handler for endpoint.
func Chain(m ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}
