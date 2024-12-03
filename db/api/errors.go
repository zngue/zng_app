package api

type Error struct {
	Code int
	Data any
	Msg  string
}
type ErrOption struct {
	Code int
	Data any
	Msg  string
}
type ErrFn func(opt *ErrOption)

func ErrCode(code int) ErrFn {
	return func(opt *ErrOption) {
		opt.Code = code
	}
}
func ErrData(data any) ErrFn {
	return func(opt *ErrOption) {
		opt.Data = data
	}
}
func ErrMsg(msg string) ErrFn {
	return func(opt *ErrOption) {
		opt.Msg = msg
	}
}
func (e *Error) Error() string {
	return e.Msg
}
func NewError(code int, data any) error {
	return &Error{
		Code: code,
		Data: data,
		Msg:  "AutoErrMsg",
	}
}
func NewErrorWithMsg(code int, data any, msg string) error {
	return &Error{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}
func NewErrFn(fns ...ErrFn) error {
	var opt = &ErrOption{
		Code: 200,
		Msg:  "AutoErrMsg",
	}
	for _, fn := range fns {
		fn(opt)
	}
	return &Error{
		Code: opt.Code,
		Data: opt.Data,
		Msg:  opt.Msg,
	}
}
