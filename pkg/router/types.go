package router

import "github.com/gin-gonic/gin"

type IRouter interface {
	Router()
}

type ApiRouterService interface {
}

type Api struct {
	router *gin.RouterGroup
	Method MethodType
	Path   string
	Fn     Fn
	IRouterServer
}
type IApiService interface {
	Register() []IRouter
}
type ApiService struct {
}
type Fn func(ctx *gin.Context) (data any, err error)
type IRouterServer struct {
}
