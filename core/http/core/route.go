package core

import "context"

type AuthMode int

const (
	AuthModeUnspecified AuthMode = iota
	AuthModeNone
	AuthModeOptional
	AuthModeRequired
)

type AuthPolicy struct {
	Mode        AuthMode
	Permissions []string
	Roles       []string
}

type RouteDesc struct {
	ServiceName   string
	MethodName    string
	Operation     string
	HTTPMethod    string
	Path          string
	Body          string
	ResponseBody  string
	Auth          AuthPolicy
	Middlewares   []string
	NewRequest    RequestFactory
	Handler       Handler
	StreamHandler StreamHandler
}

type ServiceDesc struct {
	ServiceName string
	Routes      []RouteDesc
}

type Router interface {
	GET(path string, options ...RouteOption)
	POST(path string, options ...RouteOption)
	PUT(path string, options ...RouteOption)
	DELETE(path string, options ...RouteOption)
	PATCH(path string, options ...RouteOption)
}

type ServiceRegistrar interface {
	RegisterService(serviceName string, register func(Router))
}

type RouteOption func(*RouteDesc)

func Method(methodName string) RouteOption {
	return func(route *RouteDesc) {
		route.MethodName = methodName
	}
}

func Operation(operation string) RouteOption {
	return func(route *RouteDesc) {
		route.Operation = operation
	}
}

func Body(body string) RouteOption {
	return func(route *RouteDesc) {
		route.Body = body
	}
}

func AuthNone() RouteOption {
	return func(route *RouteDesc) {
		route.Auth.Mode = AuthModeNone
	}
}

func AuthOptional() RouteOption {
	return func(route *RouteDesc) {
		route.Auth.Mode = AuthModeOptional
	}
}

func AuthRequired(permissions ...string) RouteOption {
	return func(route *RouteDesc) {
		route.Auth.Mode = AuthModeRequired
		route.Auth.Permissions = permissions
	}
}

func NewRequest[Req any]() RequestFactory {
	return func() any {
		return new(Req)
	}
}

func WrapUnary[Req any, Reply any](handler UnaryHandler[Req, Reply]) Handler {
	return func(ctx context.Context, req any) (any, error) {
		return handler(ctx, req.(*Req))
	}
}

func Unary[Req any, Reply any](handler UnaryHandler[Req, Reply]) RouteOption {
	return func(route *RouteDesc) {
		route.NewRequest = NewRequest[Req]()
		route.Handler = WrapUnary(handler)
	}
}

func Streaming[Req any](handler StreamHandlerFunc[Req]) RouteOption {
	return func(route *RouteDesc) {
		route.NewRequest = NewRequest[Req]()
		route.StreamHandler = func(ctx context.Context, req any, stream Stream) error {
			return handler(ctx, req.(*Req), stream)
		}
	}
}

type RouteCollector struct {
	serviceName string
	routes      []RouteDesc
}

func NewRouteCollector(serviceName string) *RouteCollector {
	return &RouteCollector{serviceName: serviceName}
}

func (r *RouteCollector) GET(path string, options ...RouteOption) {
	r.add("GET", path, options...)
}

func (r *RouteCollector) POST(path string, options ...RouteOption) {
	r.add("POST", path, options...)
}

func (r *RouteCollector) PUT(path string, options ...RouteOption) {
	r.add("PUT", path, options...)
}

func (r *RouteCollector) DELETE(path string, options ...RouteOption) {
	r.add("DELETE", path, options...)
}

func (r *RouteCollector) PATCH(path string, options ...RouteOption) {
	r.add("PATCH", path, options...)
}

func (r *RouteCollector) add(method string, path string, options ...RouteOption) {
	route := RouteDesc{
		ServiceName: r.serviceName,
		HTTPMethod:  method,
		Path:        path,
	}
	for _, option := range options {
		option(&route)
	}
	r.routes = append(r.routes, route)
}

func (r *RouteCollector) ServiceDesc() ServiceDesc {
	routes := make([]RouteDesc, len(r.routes))
	copy(routes, r.routes)
	return ServiceDesc{
		ServiceName: r.serviceName,
		Routes:      routes,
	}
}
