package app

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zngue/zng_app/db/api"
)

type IRouter interface {
	Router()
}
type IRouterServer struct {
}

func (IRouterServer) Router() {
	panic("Unimplemented IRouterServer.Router")
}

type RouterFn func(ctx *gin.Context) (data any, err error)

func ApiRouter(fn RouterFn) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := fn(ctx)
		if errors.Is(err, api.ErrParameter) {
			api.DataWithErr(ctx, err, data, api.Code(api.ErrorParameter))
		} else {
			api.DataWithErr(ctx, err, data)
		}
	}
}
func NewRouter(items []IApiService) (routes []IRouter) {
	for _, service := range items {
		runItems := service.Run()
		if len(runItems) > 0 {
			for _, router := range runItems {
				routes = append(routes, router)
			}
		}
	}
	return
}
