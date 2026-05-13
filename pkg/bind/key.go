package bind

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zngue/zng_app"
	"github.com/zngue/zng_app/app/server/middleware"
)

type serverContextKey struct{}        //gin.Context
type udidContextKey struct{}          //udid 参数
type localServerContextKey struct{}   //本地服务
type requestServerContextKey struct{} //请求来源服务
type operationContextKey struct{}     //请求操作

const OperationKey = "X-Request-Operation"             //请求操作
const RequestIDKey = "X-Request-Id"                    //请求ID
const LocalServer = "X-Request-Local-Server"           //当前服务名
const RequestVersion = "X-Request-Version"             //请求版本
const RequestAuthorization = "X-Request-Authorization" //请求授权
const RequestFromService = "X-Request-From-Service"    //请求来源服务

func OriginUDID() string {
	return uuid.NewString()
}

func NewServerContext(ctx context.Context, tr *gin.Context, operation string) context.Context {
	ctx = NewUDIDContext(ctx, tr)
	ctx = context.WithValue(ctx, serverContextKey{}, tr)
	ctx = context.WithValue(ctx, localServerContextKey{}, zng_app.AppName)
	ctx = context.WithValue(ctx, operationContextKey{}, operation)
	var fromService = TrHeaderFromServer(tr)
	if fromService != "" {
		ctx = context.WithValue(ctx, requestServerContextKey{}, fromService)
	}
	tr.Header(LocalServer, zng_app.AppName)
	tr.Header(OperationKey, operation)
	return ctx
}

func TrHeaderFromServer(tr *gin.Context) (str string) {
	return tr.Request.Header.Get(RequestFromService)
}

func FromRequestService(ctx context.Context) (formService string) {
	var (
		ok bool
	)
	formService, ok = ctx.Value(requestServerContextKey{}).(string)
	if ok {
		return formService
	}
	return
}

func FromServerLocalContext(ctx context.Context) (str string) {
	str, _ = ctx.Value(localServerContextKey{}).(string)
	return
}

func FromServerContext(ctx context.Context) (tr *gin.Context, ok bool) {
	tr, ok = ctx.Value(serverContextKey{}).(*gin.Context)
	return
}

func NewUDIDContext(ctx context.Context, tr *gin.Context) context.Context {
	var udid = FromUDIDContext(ctx)
	if udid == "" {
		udid = uuid.NewString()
		ctx = context.WithValue(ctx, udidContextKey{}, udid)
		tr.Header(RequestIDKey, udid)
	}
	return ctx
}

func FromUDIDContext(ctx context.Context) (str string) {
	var (
		value any
	)
	value = ctx.Value(udidContextKey{})
	str, _ = value.(string)
	return
}

func OperationByContext(ctx context.Context) (str string) {
	var (
		value any
	)
	value = ctx.Value(operationContextKey{})
	str, _ = value.(string)
	return
}
func MiddlewareHandle[R any](ctx context.Context, handler func(ctx context.Context) (R, error)) (R, error) {
	if middleware.GetRegistry() == nil {
		return handler(ctx)
	}
	operation := OperationByContext(ctx)
	wrapped := middleware.GetRegistry().Build(operation, func(ctx context.Context) (any, error) {
		return handler(ctx)
	})
	out, err := wrapped(ctx)
	if err != nil {
		var zero R
		return zero, err
	}
	return out.(R), nil
}
