package app

import "github.com/gin-gonic/gin"

type IRouter interface {
	Router()
}

type IRouterServer struct {
}

func (IRouterServer) Router() {
	panic("Unimplemented IRouterServer.Router")
}

type GinRouterFn func(ctx *gin.Context) (data any, err error)
