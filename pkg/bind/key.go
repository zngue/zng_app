package bind

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zngue/zng_app"
)

type serverGinContextKey struct{}
type udidGinContextKey struct{}
type serverNameContextKey struct{}
type operationGinContextKey struct{}

const OperationKey = "X-Request-Operation"
const RequestIDKey = "X-Request-Id"
const RequestIDServer = "X-Request-Server"
const RequestIDVersion = "X-Request-Version"
const RequestIDAuthorization = "X-Request-Authorization"

func OriginUDID() string {
	return uuid.NewString()
}

func NewServerContext(ctx context.Context, tr *gin.Context, operation string) context.Context {
	ctx = NewUDIDContext(ctx, tr)
	ctx = context.WithValue(ctx, serverGinContextKey{}, tr)               // 设置gin.Context
	ctx = context.WithValue(ctx, serverNameContextKey{}, zng_app.AppName) // 设置服务名
	ctx = context.WithValue(ctx, operationGinContextKey{}, operation)     // 设置操作
	tr.Set(OperationKey, operation)
	tr.Header(OperationKey, operation)
	tr.Header(RequestIDServer, zng_app.AppName)
	return ctx
}
func FromServerNameContext(ctx context.Context) (str string) {
	str, _ = ctx.Value(serverNameContextKey{}).(string)
	return
}
func FromServerContext(ctx context.Context) (tr *gin.Context, ok bool) {
	tr, ok = ctx.Value(serverGinContextKey{}).(*gin.Context)
	return
}

func NewUDIDContext(ctx context.Context, tr *gin.Context) context.Context {
	var udid = FromUDIDContext(ctx)
	if udid == "" {
		udid = uuid.NewString()
		ctx = context.WithValue(ctx, udidGinContextKey{}, udid)
		tr.Header(RequestIDKey, udid)
	}
	return ctx
}
func FromUDIDContext(ctx context.Context) (str string) {
	var (
		value any
	)
	value = ctx.Value(udidGinContextKey{})
	str, _ = value.(string)
	return
}

func SetOperationServerContext(tr *gin.Context, operation string) {
	tr.Set(OperationKey, operation)
}
func OperationServerContext(tr *gin.Context) (str string) {
	var (
		ok    bool
		value any
	)
	value, ok = tr.Get(OperationKey)
	if ok {
		str = value.(string)
	}
	return
}
func OperationByContext(ctx context.Context) (str string) {
	var (
		value any
	)
	value = ctx.Value(operationGinContextKey{})
	str, _ = value.(string)
	return
}
