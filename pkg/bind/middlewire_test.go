package bind

import (
	"fmt"
	"testing"

	"golang.org/x/net/context"
)

func Test_MiddleWires(t *testing.T) {
	var rs = Chain(func(h Handler) Handler {
		return func(ctx context.Context, req any) (rs any, err error) {
			fmt.Println("abc before")
			rs, err = h(ctx, req)
			fmt.Println("abc after")
			return
		}
	}, func(h Handler) Handler {
		return func(ctx context.Context, req any) (rs any, err error) {
			fmt.Println("def before")
			rs, err = h(ctx, req)
			fmt.Println("def after")
			return
		}
	}, func(h Handler) Handler {
		return func(ctx context.Context, req any) (rs any, err error) {
			fmt.Println("ghi before")
			rs, err = h(ctx, req)
			fmt.Println("ghi after")
			return
		}
	})
	var rsFn = rs(func(ctx context.Context, req any) (rs any, err error) {
		fmt.Println("jkl")
		return
	})
	rsFn(context.Background(), nil)

}
