package bind

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Bind(ctx *gin.Context, v any) (err error) {
	if ctx.Request.Method == "GET" {
		query := ctx.Request.URL.Query()
		err = binding.MapFormWithTag(v, query, "json")
		if err != nil {
			return
		}
	} else {
		err = ctx.ShouldBind(v)
	}
	return
}
