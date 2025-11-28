package bind

import (
	"context"

	"github.com/gin-gonic/gin"
)

type serverGinContextKey struct{}

func NewServerContext(ctx context.Context, tr *gin.Context) context.Context {
	return context.WithValue(ctx, serverGinContextKey{}, tr)
}

// FromServerContext returns the Transport value stored in ctx, if any.
func FromServerContext(ctx context.Context) (tr *gin.Context, ok bool) {
	tr, ok = ctx.Value(serverGinContextKey{}).(*gin.Context)
	return
}
