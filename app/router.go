package app

import (
	"github.com/gin-gonic/gin"
	"github.com/zngue/zng_app/db/api"
	"github.com/zngue/zng_app/errors"
	"github.com/zngue/zng_app/log"
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
		if err != nil {
			if SyncLogger { // api接口错误自动记录日志
				//异步记录日志
				rsErr := errors.Wrap(err, "api data err")
				go func() {
					defer func() {
						if r := recover(); r != nil {
							log.Errorf("panic recover err:%v", r)
						}
					}()
					log.Error(errors.GetStackTrace(rsErr))
				}()
			}
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
