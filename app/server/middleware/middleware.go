package middleware

import (
	"context"
)

type Handler func(ctx context.Context) (any, error)

type RequestHandler[R any] func(ctx context.Context) (R, error)

type Middleware func(ctx context.Context, next Handler) (any, error)

type AfterHandler func(ctx context.Context, err error, in any, rs any)

type AfterMiddleware func(ctx context.Context, err error, in any, rs any, next AfterHandler)

var globalRegistry *Registry

func SetRegistry(r *Registry) {
	globalRegistry = r
}

func GetRegistry() *Registry {
	return globalRegistry
}

type Registry struct {
	before      []Middleware
	after       []AfterMiddleware
	beforePerOp map[string][]Middleware
	afterPerOp  map[string][]AfterMiddleware
}

type RegistryOption func(*Registry)

func WithBefore(mw ...Middleware) RegistryOption {
	return func(r *Registry) {
		r.before = append(r.before, mw...)
	}
}

func WithBeforeOperation(op string, mw ...Middleware) RegistryOption {
	return func(r *Registry) {
		if r.beforePerOp == nil {
			r.beforePerOp = make(map[string][]Middleware)
		}
		r.beforePerOp[op] = append(r.beforePerOp[op], mw...)
	}
}

func WithAfter(mw ...AfterMiddleware) RegistryOption {
	return func(r *Registry) {
		r.after = append(r.after, mw...)
	}
}

func WithAfterOperation(op string, mw ...AfterMiddleware) RegistryOption {
	return func(r *Registry) {
		if r.afterPerOp == nil {
			r.afterPerOp = make(map[string][]AfterMiddleware)
		}
		r.afterPerOp[op] = append(r.afterPerOp[op], mw...)
	}
}

func NewRegistry(opts ...RegistryOption) *Registry {
	r := &Registry{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *Registry) Build(operation string, handler Handler) Handler {
	h := handler
	if ops, ok := r.beforePerOp[operation]; ok {
		for i := len(ops) - 1; i >= 0; i-- {
			h = wrap(ops[i], h)
		}
	}
	for i := len(r.before) - 1; i >= 0; i-- {
		h = wrap(r.before[i], h)
	}
	return h
}

func (r *Registry) BuildAfter(operation string, handler AfterHandler) AfterHandler {
	h := handler
	if ops, ok := r.afterPerOp[operation]; ok {
		for i := len(ops) - 1; i >= 0; i-- {
			h = wrapAfter(ops[i], h)
		}
	}
	for i := len(r.after) - 1; i >= 0; i-- {
		h = wrapAfter(r.after[i], h)
	}
	return h
}

func wrap(mw Middleware, next Handler) Handler {
	return func(ctx context.Context) (any, error) {
		return mw(ctx, next)
	}
}

func wrapAfter(mw AfterMiddleware, next AfterHandler) AfterHandler {
	return func(ctx context.Context, err error, in any, rs any) {
		mw(ctx, err, in, rs, next)
	}
}

func SkipWhen(condition func(ctx context.Context) bool, mw Middleware) Middleware {
	return func(ctx context.Context, next Handler) (any, error) {
		if condition(ctx) {
			return next(ctx)
		}
		return mw(ctx, next)
	}
}

func SkipAfterWhen(condition func(ctx context.Context) bool, mw AfterMiddleware) AfterMiddleware {
	return func(ctx context.Context, err error, in any, rs any, next AfterHandler) {
		if condition(ctx) {
			next(ctx, err, in, rs)
			return
		}
		mw(ctx, err, in, rs, next)
	}
}
