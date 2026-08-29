package core

import "context"

type operationContextKey struct{}
type requestIDContextKey struct{}

const RequestIDHeader = "X-Request-Id"

type OperationInfo struct {
	ServiceName string
	MethodName  string
	Operation   string
	HTTPMethod  string
	Path        string
}

func WithOperation(ctx context.Context, route RouteDesc) context.Context {
	return context.WithValue(ctx, operationContextKey{}, OperationInfo{
		ServiceName: route.ServiceName,
		MethodName:  route.MethodName,
		Operation:   route.Operation,
		HTTPMethod:  route.HTTPMethod,
		Path:        route.Path,
	})
}

func OperationFromContext(ctx context.Context) (OperationInfo, bool) {
	op, ok := ctx.Value(operationContextKey{}).(OperationInfo)
	return op, ok
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDContextKey{}, requestID)
}

func RequestIDFromContext(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDContextKey{}).(string)
	return requestID
}
