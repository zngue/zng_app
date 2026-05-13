package bind

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code      MessageCode `json:"statusCode"`          // 业务状态码（0=成功）
	Message   string      `json:"message"`             // 提示信息
	Reason    string      `json:"reason,omitempty"`    // 错误原因
	RequestId string      `json:"requestId,omitempty"` // 请求追踪ID
	Data      any         `json:"data,omitempty"`      // 实际业务数据
}

type MessageCode int

// 预定义业务状态码
const (
	CodeSuccess       MessageCode = 200   // 成功
	ErrorParameter    MessageCode = 422   // 参数错误
	ErrorResponse     MessageCode = 400   //响应错误
	CodeParamError    MessageCode = 10001 // 参数错误
	CodeUnauthorized  MessageCode = 10002 // 未授权
	CodeForbidden     MessageCode = 10003 // 无权限
	CodeNotFound      MessageCode = 10004 // 资源不存在
	CodeInternalError MessageCode = 10500 // 内部错误

)

type ResOption struct {
	Code    MessageCode
	Message string
	Reason  string
	Data    any
}
type ResOptionFn func(opt *ResOption)

func Data(data any) ResOptionFn {
	return func(opt *ResOption) {
		opt.Data = data
	}
}
func DataCode(code MessageCode) ResOptionFn {
	return func(opt *ResOption) {
		opt.Code = code
	}
}
func DataMsg(msg string) ResOptionFn {
	return func(opt *ResOption) {
		opt.Message = msg
	}
}
func DataReason(err error) ResOptionFn {
	return func(opt *ResOption) {
		if err != nil {
			opt.Reason = err.Error()
		}
	}
}
func (m MessageCode) String() string {
	switch m {
	case CodeSuccess:
		return "success"
	case ErrorParameter:
		return "参数错误"
	case ErrorResponse:
		return "响应错误"
	case CodeParamError:
		return "参数错误"
	case CodeUnauthorized:
		return "未授权"
	case CodeForbidden:
		return "无权限"
	case CodeNotFound:
		return "资源不存在"
	case CodeInternalError:
		return "内部错误"
	default:
		return "未知错误"
	}
}

func ApiDataWithErr[T any](ctx *gin.Context, err error, data T, fns ...ResOptionFn) {
	if err != nil {
		ApiErrorResponse(ctx, err, fns...)
	} else {
		ApiCodeSuccess(ctx, data, fns...)
	}
}
func ApiErrorResponse(ctx *gin.Context, err error, fns ...ResOptionFn) {
	var res = &ResOption{
		Code: CodeInternalError,
	}
	fns = append(fns, DataCode(ErrorResponse))
	fns = append(fns, DataReason(err))
	for _, fn := range fns {
		fn(res)
	}
	Execute(ctx, res)

}
func ApiCodeInternalError(ctx *gin.Context, err error, fns ...ResOptionFn) {
	var res = &ResOption{
		Code: CodeInternalError,
	}
	res.Code = CodeInternalError
	for _, fn := range fns {
		fn(res)
	}
	Execute(ctx, res)
}

func ApiCodeSuccess(ctx *gin.Context, data any, fns ...ResOptionFn) {
	var res = &ResOption{
		Code: CodeSuccess,
	}
	if data != nil {
		fns = append(fns, Data(data))
	}
	for _, fn := range fns {
		fn(res)
	}
	Execute(ctx, res)
}

func ApiErrorParameter(ctx *gin.Context, err error, fns ...ResOptionFn) {
	var res = &ResOption{
		Code: CodeSuccess,
	}
	res.Code = CodeParamError
	if err != nil {
		fns = append(fns, DataReason(err))
	}
	for _, fn := range fns {
		fn(res)
	}
	Execute(ctx, res)
}

func Execute(ctx *gin.Context, in *ResOption) {
	var requestId string
	if ctx.Writer != nil {
		requestId = ctx.Writer.Header().Get(RequestIDKey)
	}
	var out = &Response{
		Code:      in.Code,
		Message:   in.Message,
		Reason:    in.Reason,
		RequestId: requestId,
		Data:      in.Data,
	}
	if out.Message == "" {
		out.Message = out.Code.String()
	}
	ctx.JSON(http.StatusOK, out)
}

// WeChatPayError /*
func WeChatPayError(ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"code":    "FAILED",
		"message": "支付失败",
	})
}
func WeChatPaySuccess(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "成功",
	})
}
