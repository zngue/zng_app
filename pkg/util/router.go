package util

import (
	"github.com/gin-gonic/gin"
	"github.com/zngue/zng_app/pkg/router"
)

func RouterGroupFn[T any](path string, api *gin.RouterGroup) router.GroupRouter[T] {
	return api.Group(path)
}
