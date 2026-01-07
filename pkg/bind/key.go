package bind

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zngue/zng_app"
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
	ctx = context.WithValue(ctx, serverContextKey{}, tr)                   // 设置gin.Context
	ctx = context.WithValue(ctx, localServerContextKey{}, zng_app.AppName) // 设置服务名
	ctx = context.WithValue(ctx, operationContextKey{}, operation)         // 设置操作
	var fromService = TrHeaderFromServer(tr)
	if fromService != "" {
		ctx = context.WithValue(ctx, requestServerContextKey{}, fromService)
	}
	tr.Set(OperationKey, operation)         // 设置操作
	tr.Header(OperationKey, operation)      // 设置当前操作在头部
	tr.Header(LocalServer, zng_app.AppName) // 设置当前服务在头部
	return ctx
}

// 获取头部请求来源
func TrHeaderFromServer(tr *gin.Context) (str string) {
	return tr.Request.Header.Get(RequestFromService)
}

// 获取来源服务
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

func FromServerLocalContext(ctx context.Context) (str string) { //获取当前服务名
	str, _ = ctx.Value(localServerContextKey{}).(string)
	return
}
func FromServerContext(ctx context.Context) (tr *gin.Context, ok bool) {
	tr, ok = ctx.Value(localServerContextKey{}).(*gin.Context)
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
