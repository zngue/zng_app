package bind

import (
	"context"
	"time"

	"github.com/zngue/zng_app/app/server/middleware"
	"github.com/zngue/zng_app/log"
	"go.uber.org/zap"
)

func LogAfterMiddleware(ctx context.Context, err error, in any, rs any, next middleware.AfterHandler) {
	next(ctx, err, in, rs)
	operation := OperationByContext(ctx)
	path := PathByContext(ctx)
	requestId := FromUDIDContext(ctx)
	start := StartTimeByContext(ctx)
	method := "GET"
	var duration time.Duration
	if !start.IsZero() {
		duration = time.Since(start)
	}
	if tr, ok := FromServerContext(ctx); ok {
		method = tr.Request.Method
	}
	fields := []zap.Field{
		zap.String("method", method),
		zap.String("operation", operation),
		zap.String("path", path),
		zap.String("requestId", requestId),
		zap.Duration("duration", duration),
		zap.Any("params", in),
	}
	if err != nil {
		fields = append(fields, zap.Error(err))
		log.Default().Error("request_error", fields...)
		return
	}
	log.Default().Info("request_success", fields...)
}
