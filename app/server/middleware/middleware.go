package middleware

import "context"

type Handler func(ctx context.Context) (any, error)

type Middleware func(ctx context.Context, next Handler) (any, error)

var globalRegistry *Registry

func SetRegistry(r *Registry) {
	globalRegistry = r
}

func GetRegistry() *Registry {
	return globalRegistry
}

type Registry struct {
	global []Middleware
	perOp  map[string][]Middleware
}

type RegistryOption func(*Registry)

func WithGlobal(mw ...Middleware) RegistryOption {
	return func(r *Registry) {
		r.global = append(r.global, mw...)
	}
}

func WithOperation(op string, mw ...Middleware) RegistryOption {
	return func(r *Registry) {
		if r.perOp == nil {
			r.perOp = make(map[string][]Middleware)
		}
		r.perOp[op] = append(r.perOp[op], mw...)
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
	if ops, ok := r.perOp[operation]; ok {
		for i := len(ops) - 1; i >= 0; i-- {
			h = wrap(ops[i], h)
		}
	}
	for i := len(r.global) - 1; i >= 0; i-- {
		h = wrap(r.global[i], h)
	}
	return h
}

func wrap(mw Middleware, next Handler) Handler {
	return func(ctx context.Context) (any, error) {
		return mw(ctx, next)
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

type MiddlewareChain struct {
	middlewares []Middleware
}

func NewChain(fns ...Middleware) *MiddlewareChain {
	return &MiddlewareChain{middlewares: fns}
}

func (c *MiddlewareChain) Handle(ctx context.Context, handler Handler) (any, error) {
	h := handler
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		h = wrap(c.middlewares[i], h)
	}
	return h(ctx)
}
