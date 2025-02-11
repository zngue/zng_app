package bind

import "context"

type MiddlewareFn = func(ctx context.Context) (tx context.Context, err error)

var middlewares []MiddlewareFn

func SetMiddleWires(fns ...MiddlewareFn) {
	middlewares = fns
	return
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
