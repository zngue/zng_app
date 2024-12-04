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
	if r.Method == "" {
		r.router.Any(r.Path, ApiRouter(r.Fn))
	} else {
		r.router.Handle(string(r.Method), r.Path, ApiRouter(r.Fn))
	}
}
func ApiFn(router *gin.RouterGroup, method MethodType, path string, fn RouterFn) *Api {
	return &Api{router: router, Method: method, Path: path, Fn: fn}
}
func ApiAnyFn(router *gin.RouterGroup, path string, fn RouterFn) *Api {
	return &Api{router: router, Path: path, Fn: fn}
}
func ApiGetFn(router *gin.RouterGroup, path string, fn RouterFn) *Api {
	return &Api{router: router, Method: GET, Path: path, Fn: fn}
}
func ApiPostFn(router *gin.RouterGroup, path string, fn RouterFn) *Api {
	return &Api{router: router, Method: POST, Path: path, Fn: fn}
}
