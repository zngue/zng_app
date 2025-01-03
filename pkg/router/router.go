package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zngue/zng_app/db/api"
	"github.com/zngue/zng_app/errors"
	"net/http"
)

func (IRouterServer) Router() {
	panic("Unimplemented IRouterServer.Router")
}

type MethodType string

const (
	POST MethodType = http.MethodPost
	GET  MethodType = http.MethodGet
)

func (a *ApiService) Register() []*Api {
	panic("implement me")
}
func ApiServiceFn(dataItems ...*Api) []*Api {
	return dataItems
}

func (r *Api) Router() {
	if r.Method == "" {
		r.router.Any(r.Path, ApiRouter(r.Fn))
	} else {
		r.router.Handle(string(r.Method), r.Path, ApiRouter(r.Fn))
	}
}
func ApiFn(router *gin.RouterGroup, method MethodType, path string, fn Fn) *Api {
	return &Api{router: router, Method: method, Path: path, Fn: fn}
}
func ApiAnyFn(router *gin.RouterGroup, path string, fn Fn) *Api {
	return &Api{router: router, Path: path, Fn: fn}
}
func ApiGetFn(router *gin.RouterGroup, path string, fn Fn) *Api {
	return &Api{router: router, Method: GET, Path: path, Fn: fn}
}
func ApiPostFn(router *gin.RouterGroup, path string, fn Fn) *Api {
	return &Api{router: router, Method: POST, Path: path, Fn: fn}
}

func ApiRouter(fn Fn) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := fn(ctx)
		if err != nil {
			errors.LogS(err)
		}
		if errors.Is(err, api.ErrParameter) {
			api.DataApiWithErr(ctx, err, data, api.Code(api.ErrorParameter))
			return
		} else {
			api.DataApiWithErr(ctx, err, data)
			return
		}
	}
}
