package bind

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/zngue/zng_app/errors"
	"golang.org/x/net/context"
)

func Bind(ctx *gin.Context, v any) (err error) {
	if ctx.Request.Method == http.MethodGet {
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

func GinContext(ctx context.Context) (ginCtx *gin.Context, err error) {
	value := ctx.Value("gin_ctx")
	var ok bool
	ginCtx, ok = value.(*gin.Context)
	if ok {
		return
	}
	err = errors.New("gin.Context is not exist")
	return
}

// GetOperation 获取操作对象
func GetOperation(ctx context.Context) (operation string, err error) {
	value := ctx.Value("operation")
	var ok bool
	operation, ok = value.(string)
	if ok {
		return
	}
	err = errors.New("operation is not exist")
	return
}
