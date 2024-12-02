package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type MethodType string

const (
	POST MethodType = http.MethodPost
	GET  MethodType = http.MethodGet
)

type IApiService interface {
	Run() []*Api
}
type ApiService struct {
}

func (a *ApiService) Run() []*Api {
	panic("implement me")
}
func ApiServiceFn(dataItems ...*Api) []*Api {
	return dataItems
}

type Api struct {
	router *gin.RouterGroup
	Method MethodType
	Path   string
	Fn     RouterFn
	IRouterServer
}

func (r *Api) Router() {
	r.router.Handle(string(r.Method), r.Path, ApiRouter(r.Fn))
}
func ApiFn(router *gin.RouterGroup, method MethodType, path string, fn RouterFn) *Api {
	return &Api{router: router, Method: method, Path: path, Fn: fn}
}
