package core

import "context"

import "github.com/zngue/zng_app/core/errors_ez"

type TransportContext interface {
	Context() context.Context
	SetContext(context.Context)
	Param(name string) string
	Query(name string) string
	QueryArray(name string) []string
	Header(name string) string
	SetHeader(name string, value string)
	BindJSON(obj any) error
	JSON(statusCode int, obj any)
	AbortWithJSON(statusCode int, obj any)
	NewStream(ctx context.Context) Stream
}

type Binder interface {
	Bind(ctx TransportContext, route RouteDesc, req any) error
}

type Responder interface {
	Success(ctx TransportContext, data any)
	Error(ctx TransportContext, err error)
}

type Authenticator interface {
	Authenticate(ctx TransportContext, route RouteDesc) (context.Context, error)
}

type RequestFactory func() any

type Handler func(ctx context.Context, req any) (any, error)

type UnaryHandler[Req any, Reply any] func(ctx context.Context, req *Req) (*Reply, error)

type Middleware func(ctx context.Context, route RouteDesc, req any, next Handler) (any, error)

type Response struct {
	StatusCode int                    `json:"statusCode"`
	Reason     string                 `json:"reason,omitempty"`
	Message    string                 `json:"message"`
	RequestID  string                 `json:"requestId,omitempty"`
	Data       any                    `json:"data,omitempty"`
	Stack      []errors_ez.StackFrame `json:"stack,omitempty"`
}

type StackFrame = errors_ez.StackFrame
