package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/zngue/zng_app/core/errors_ez"
	"github.com/zngue/zng_app/core/http/core"
	"github.com/zngue/zng_app/core/http/validate"
)

var requestIDCounter atomic.Uint64

func NewRequestID() string {
	seq := requestIDCounter.Add(1)
	return fmt.Sprintf("req-%d-%d", time.Now().UnixNano(), seq)
}

type Runtime struct {
	binder        core.Binder
	validator     validate.Validator
	responder     core.Responder
	authenticator core.Authenticator
	middlewares   []core.Middleware
}

type Option func(*Runtime)

func NewRuntime(options ...Option) *Runtime {
	rt := &Runtime{
		binder:        DefaultBinder{},
		validator:     validate.NoopValidator{},
		responder:     JSONResponder{},
		authenticator: NoopAuthenticator{},
	}
	for _, option := range options {
		option(rt)
	}
	return rt
}

func WithBinder(binder core.Binder) Option {
	return func(rt *Runtime) {
		rt.binder = binder
	}
}

func WithValidator(validator validate.Validator) Option {
	return func(rt *Runtime) {
		rt.validator = validator
	}
}

func WithResponder(responder core.Responder) Option {
	return func(rt *Runtime) {
		rt.responder = responder
	}
}

func WithDebug(debug bool) Option {
	return func(rt *Runtime) {
		if r, ok := rt.responder.(JSONResponder); ok {
			r.Debug = debug
			rt.responder = r
		}
	}
}

func WithAuthenticator(authenticator core.Authenticator) Option {
	return func(rt *Runtime) {
		rt.authenticator = authenticator
	}
}

func WithMiddleware(middlewares ...core.Middleware) Option {
	return func(rt *Runtime) {
		rt.middlewares = append(rt.middlewares, middlewares...)
	}
}

func (rt *Runtime) Handle(ctx core.TransportContext, route core.RouteDesc) {
	if route.NewRequest == nil {
		rt.responder.Error(ctx, errors_ez.Internal("route %s missing request constructor", route.Operation))
		return
	}

	if route.StreamHandler != nil {
		rt.handleStream(ctx, route)
		return
	}

	if route.Handler == nil {
		rt.responder.Error(ctx, errors_ez.Internal("route %s missing handler", route.Operation))
		return
	}

	requestID := strings.TrimSpace(ctx.Header(core.RequestIDHeader))
	if requestID == "" {
		requestID = NewRequestID()
	}
	ctx.SetHeader(core.RequestIDHeader, requestID)

	req := route.NewRequest()

	if err := rt.binder.Bind(ctx, route, req); err != nil {
		rt.responder.Error(ctx, errors_ez.FromError(err, http.StatusBadRequest, "BIND_FAILED"))
		return
	}

	if err := rt.validator.Validate(req); err != nil {
		rt.responder.Error(ctx, errors_ez.BadRequest("VALIDATE_FAILED", "%s", err))
		return
	}

	baseCtx := core.WithRequestID(core.WithOperation(ctx.Context(), route), requestID)
	ctx.SetContext(baseCtx)
	baseCtx, err := rt.authenticator.Authenticate(ctx, route)
	if err != nil {
		rt.responder.Error(ctx, err)
		return
	}
	baseCtx = core.WithRequestID(baseCtx, requestID)
	ctx.SetContext(baseCtx)

	handler := route.Handler

	for i := len(rt.middlewares) - 1; i >= 0; i-- {
		middleware := rt.middlewares[i]
		next := handler
		handler = func(callCtx context.Context, callReq any) (any, error) {
			return middleware(callCtx, route, callReq, next)
		}
	}

	data, err := handler(baseCtx, req)
	if err != nil {
		rt.responder.Error(ctx, err)
		return
	}
	rt.responder.Success(ctx, data)
}

func (rt *Runtime) handleStream(ctx core.TransportContext, route core.RouteDesc) {
	requestID := strings.TrimSpace(ctx.Header(core.RequestIDHeader))
	if requestID == "" {
		requestID = NewRequestID()
	}
	ctx.SetHeader(core.RequestIDHeader, requestID)
	ctx.SetHeader("Content-Type", "text/event-stream")
	ctx.SetHeader("Cache-Control", "no-cache")
	ctx.SetHeader("Connection", "keep-alive")

	req := route.NewRequest()

	if err := rt.binder.Bind(ctx, route, req); err != nil {
		rt.responder.Error(ctx, errors_ez.FromError(err, http.StatusBadRequest, "BIND_FAILED"))
		return
	}

	if err := rt.validator.Validate(req); err != nil {
		rt.responder.Error(ctx, errors_ez.BadRequest("VALIDATE_FAILED", "%s", err))
		return
	}

	baseCtx := core.WithRequestID(core.WithOperation(ctx.Context(), route), requestID)
	ctx.SetContext(baseCtx)
	baseCtx, err := rt.authenticator.Authenticate(ctx, route)
	if err != nil {
		rt.responder.Error(ctx, err)
		return
	}
	baseCtx = core.WithRequestID(baseCtx, requestID)
	ctx.SetContext(baseCtx)

	stream := ctx.NewStream(baseCtx)
	if err := route.StreamHandler(baseCtx, req, stream); err != nil {
		_ = stream.Send(errors_ez.ExtractError(err))
	}
	_ = stream.Close()
}

type DefaultBinder struct{}

func (DefaultBinder) Bind(ctx core.TransportContext, route core.RouteDesc, req any) error {
	if route.Body != "" {
		return ctx.BindJSON(req)
	}
	return core.BindQuery(ctx, req)
}

type NoopAuthenticator struct{}

func (NoopAuthenticator) Authenticate(ctx core.TransportContext, route core.RouteDesc) (context.Context, error) {
	return ctx.Context(), nil
}

type JSONResponder struct {
	Debug bool
}

func (r JSONResponder) Success(ctx core.TransportContext, data any) {
	ctx.JSON(http.StatusOK, core.Response{
		StatusCode: http.StatusOK,
		Message:    "success",
		RequestID:  core.RequestIDFromContext(ctx.Context()),
		Data:       data,
	})
}

func (r JSONResponder) Error(ctx core.TransportContext, err error) {
	apiErr := errors_ez.ExtractError(err)

	stack := apiErr.Stack
	message := apiErr.Message

	if !r.Debug {
		stack = nil
		if apiErr.Code >= 500 {
			message = "internal server error"
		}
	}

	ctx.AbortWithJSON(apiErr.Code, core.Response{
		StatusCode: apiErr.Code,
		Reason:     apiErr.Reason,
		Message:    message,
		Stack:      stack,
		RequestID:  core.RequestIDFromContext(ctx.Context()),
	})
}
