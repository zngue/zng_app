package middleware

import "context"

type Middleware func(ctx context.Context) (context.Context, error)

type MiddlewareChain struct {
	middlewares []Middleware
}

func NewChain(fns ...Middleware) *MiddlewareChain {
	return &MiddlewareChain{middlewares: fns}
}

func SkipWhen(condition func(ctx context.Context) bool, f Middleware) Middleware {
	return func(ctx context.Context) (context.Context, error) {
		if condition(ctx) {
			return ctx, nil
		}
		return f(ctx)
	}
}

func (c *MiddlewareChain) Handle(ctx context.Context) (context.Context, error) {
	for _, f := range c.middlewares {
		ctx, err := f(ctx)
		if err != nil {
			return ctx, err
		}
	}
	return ctx, nil
}
