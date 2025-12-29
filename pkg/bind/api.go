package bind

import (
	"github.com/gin-gonic/gin"
	"github.com/zngue/zng_app/errors"
)

var (
	SuccessCode    = 200
	ErrorCode      = 100
	ErrorParameter = 422
	SuccessMsg     = "success"
	ErrorMsg       = "error"
	ParameterMsg   = "参数错误"
	ErrParameter   = errors.New("请求参数错误")
)

type Response struct {
	Code int    `json:"statusCode" `
	Msg  string `json:"message" `
	Data any    `json:"data" `
}

type Fn func(response *Response)

func Code(code int) Fn {
	return func(response *Response) {
		response.Code = code
	}
}
func Err(err error) Fn {
	return func(response *Response) {
		if err != nil {
			response.Msg = err.Error()
		}
	}
}
func Data(data any) Fn {
	return func(response *Response) {
		response.Data = data
	}
}

// Msg /*
func Msg(msg string) Fn {
	return func(response *Response) {
		response.Msg = msg
	}
}
func DataError(ctx *gin.Context, err error, fns ...Fn) {
	var data = &Response{
		Code: ErrorCode,
		Msg:  ErrorMsg,
		Data: nil,
	}
	if err != nil {
		fns = append(fns, Err(err))
	}
	if len(fns) > 0 {
		for _, fn := range fns {
			fn(data)
		}
	}
	ctx.JSON(200, data)
}
func DataApiWithErr(ctx *gin.Context, err error, data any, fns ...Fn) {
	if err != nil {
		DataError(ctx, err, fns...)
	} else {
		if data == nil {
			DataSuccess(ctx, &Response{
				Code: SuccessCode,
				Msg:  SuccessMsg,
			})
		} else {
			DataSuccess(ctx, data)
		}
	}
}
func DataSuccess(ctx *gin.Context, data any) {
	ctx.JSON(200, data)
}
