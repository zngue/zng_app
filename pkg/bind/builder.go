package bind

import (
	"github.com/zngue/go_helper/pkg/util"
	"golang.org/x/net/context"
)

func Server(ms ...Middleware) *Builder {
	return &Builder{middlewares: ms}
}

type Builder struct {
	middlewares []Middleware
	path        []string
}

func (b *Builder) Path(paths ...string) {
	b.path = paths
}
func (b *Builder) Build() Middleware {
	return Selector(b.path, b.middlewares...)
}
func Selector(paths []string, ms ...Middleware) Middleware {
	return func(handler Handler) Handler {
		return func(ctx context.Context, req any) (rs any, err error) {
			var operation = OperationByContext(ctx)
			if len(paths) > 0 && operation != "" && util.InArray(operation, paths) {
				return handler(ctx, req)
			}
			if len(ms) > 0 {
				return Chain(ms...)(handler)(ctx, req)
			}
			return handler(ctx, req)
		}
	}
}
