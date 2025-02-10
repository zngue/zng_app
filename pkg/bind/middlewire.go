package bind

import "context"

type MiddlewareFn = func(ctx context.Context) (err error)

var middlewares []MiddlewareFn

func SetMiddleWires(fns ...MiddlewareFn) {
	middlewares = fns
	return
}
func GetMiddleWires() []MiddlewareFn {
	return middlewares
}
func ClearMiddleWires() {
	middlewares = nil
}
