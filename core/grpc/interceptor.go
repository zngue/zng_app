package grpc

import (
	"context"
	"log"
	"time"

	"github.com/zngue/zng_app/core/errors_ez"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	resp, err := handler(ctx, req)
	if err == nil {
		return resp, nil
	}

	apiErr := errors_ez.ExtractError(err)

	if apiErr.Reason != "" {
		return nil, status.Error(codeFromHTTP(apiErr.Code), apiErr.Reason)
	}
	return nil, status.Error(codeFromHTTP(apiErr.Code), apiErr.Message)
}

func LoggingInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("grpc method=%s cost=%s err=%v", info.FullMethod, time.Since(start), err)
	return resp, err
}

func ChainUnary(opts ...grpc.ServerOption) []grpc.ServerOption {
	return append([]grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			LoggingInterceptor,
			ErrorInterceptor,
		),
	}, opts...)
}

func codeFromHTTP(httpCode int) codes.Code {
	switch {
	case httpCode >= 200 && httpCode < 300:
		return codes.OK
	case httpCode == 400:
		return codes.InvalidArgument
	case httpCode == 401:
		return codes.Unauthenticated
	case httpCode == 403:
		return codes.PermissionDenied
	case httpCode == 404:
		return codes.NotFound
	case httpCode == 409:
		return codes.AlreadyExists
	case httpCode == 429:
		return codes.ResourceExhausted
	case httpCode >= 500:
		return codes.Internal
	default:
		return codes.Unknown
	}
}
