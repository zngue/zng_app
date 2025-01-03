package router

import "github.com/gin-gonic/gin"

type IRouter interface {
	Router()
}

type Api struct {
	router *gin.RouterGroup
	Method MethodType
	Path   string
	Fn     Fn
	IRouterServer
}
type IApiService interface {
	Register() []*Api
}
type ApiService struct {
}
type Fn func(ctx *gin.Context) (data any, err error)
type IRouterServer struct {
}
type GroupRouter[T any] *gin.RouterGroup

func GroupRouterFn[T any](path string, api *gin.RouterGroup) GroupRouter[T] {
	return api.Group(path)
}
