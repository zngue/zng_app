package errors_ez

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/zngue/zng_app"
)

type StackFrame struct {
	File     string `json:"file,omitempty"`
	Line     int    `json:"line,omitempty"`
	Function string `json:"function,omitempty"`
	Message  string `json:"message,omitempty"`
}

type StackTracer interface {
	LastMessage() string
	StackTrace() []StackFrame
}

type Error struct {
	Code    int          `json:"code"`
	Reason  string       `json:"reason,omitempty"`
	Message string       `json:"message"`
	Stack   []StackFrame `json:"stack,omitempty"`
	cause   error
}

func (e *Error) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("%s: %s", e.Reason, e.Message)
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.cause
}

func (e *Error) WithCode(code int, reason string) *Error {
	e.Code = code
	e.Reason = reason
	return e
}

func newError(code int, reason, format string, args ...any) *Error {
	return &Error{
		Code:    code,
		Reason:  reason,
		Message: fmt.Sprintf(format, args...),
		Stack:   CaptureStack(3),
	}
}

func BadRequest(reason, format string, args ...any) *Error {
	return newError(http.StatusBadRequest, reason, format, args...)
}

func Unauthorized(reason, format string, args ...any) *Error {
	return newError(http.StatusUnauthorized, reason, format, args...)
}

func Forbidden(reason, format string, args ...any) *Error {
	return newError(http.StatusForbidden, reason, format, args...)
}

func NotFound(reason, format string, args ...any) *Error {
	return newError(http.StatusNotFound, reason, format, args...)
}

func Conflict(reason, format string, args ...any) *Error {
	return newError(http.StatusConflict, reason, format, args...)
}

func Internal(format string, args ...any) *Error {
	return newError(http.StatusInternalServerError, "INTERNAL", format, args...)
}

func New(format string, args ...any) *Error {
	msg := fmt.Sprintf(format, args...)
	return &Error{
		Code:    http.StatusInternalServerError,
		Reason:  "INTERNAL",
		Message: msg,
		Stack:   CaptureStack(2),
	}
}

func Wrap(err error) *Error {
	return wrapErr(err, "", false)
}

func WrapF(err error, format string, args ...any) *Error {
	return wrapErr(err, fmt.Sprintf(format, args...), true)
}

func wrapErr(err error, msg string, hasMsg bool) *Error {
	if err == nil {
		if hasMsg {
			return New("%s", msg)
		}
		return nil
	}

	var existing *Error
	if errors.As(err, &existing) {
		frame := CaptureStack(3)[0]
		if hasMsg {
			frame.Message = msg
		}
		existing.Stack = append(existing.Stack, frame)
		if hasMsg {
			existing.Message = msg
		}
		return existing
	}

	var stack []StackFrame
	var tracer StackTracer
	if errors.As(err, &tracer) {
		stack = tracer.StackTrace()
	}

	frame := CaptureStack(3)[0]
	if hasMsg {
		frame.Message = msg
	}
	stack = append(stack, frame)

	message := err.Error()
	if hasMsg {
		message = msg
	}

	return &Error{
		Code:    http.StatusInternalServerError,
		Reason:  "INTERNAL",
		Message: message,
		Stack:   stack,
		cause:   err,
	}
}

func FromError(err error, code int, reason string) *Error {
	var apiErr *Error
	if errors.As(err, &apiErr) {
		return apiErr
	}

	var tracer StackTracer
	if errors.As(err, &tracer) {
		msg := tracer.LastMessage()
		if msg == "" {
			if unwrapped := errors.Unwrap(err); unwrapped != nil {
				msg = unwrapped.Error()
			} else {
				msg = err.Error()
			}
		}
		return &Error{
			Code:    code,
			Reason:  reason,
			Message: msg,
			Stack:   tracer.StackTrace(),
		}
	}

	return &Error{
		Code:    code,
		Reason:  reason,
		Message: err.Error(),
		Stack:   CaptureStack(2),
	}
}

func CaptureStack(skip int) []StackFrame {
	var frames []StackFrame
	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		funcName := ""
		if fn != nil {
			funcName = fn.Name()
		}
		frames = append(frames, StackFrame{
			File:     file,
			Line:     line,
			Function: funcName,
		})
	}
	return frames
}

func ExtractError(err error) *Error {
	var apiErr *Error
	if errors.As(err, &apiErr) {
		return apiErr
	}

	var tracer StackTracer
	if errors.As(err, &tracer) {
		msg := tracer.LastMessage()
		if msg == "" {
			if unwrapped := errors.Unwrap(err); unwrapped != nil {
				msg = unwrapped.Error()
			} else {
				msg = err.Error()
			}
		}
		return &Error{
			Code:    http.StatusInternalServerError,
			Reason:  "INTERNAL",
			Message: msg,
			Stack:   tracer.StackTrace(),
		}
	}

	return &Error{
		Code:    http.StatusInternalServerError,
		Reason:  "INTERNAL",
		Message: err.Error(),
	}
}

type EzError struct {
	customStr []string
	index     int
	lineErr   []*EzInfo
}

type EzInfo struct {
	ErrInfo string `json:"errInfo"`
	FnInfo  string `json:"fnName"`
	Message string `json:"message,omitempty"`
	Index   int    `json:"index"`
	AppName string `json:"appName,omitempty"`
}

func (e *EzError) ReasonMessage() []*EzInfo {
	return e.lineErr
}

func (e *EzError) Error() string {
	byt, err := json.Marshal(e.lineErr)
	if err != nil {
		return fmt.Sprintf("%+v", e.lineErr)
	}
	return string(byt)
}

func (e *EzError) LastCustomReason() string {
	if len(e.customStr) > 0 {
		return e.customStr[len(e.customStr)-1]
	}
	return ""
}

func (e *EzError) Unwrap() error {
	if len(e.lineErr) > 0 && e.lineErr[0].Message != "" {
		return fmt.Errorf("%s", e.lineErr[0].Message)
	}
	return nil
}

func (e *EzError) LastMessage() string {
	return e.LastCustomReason()
}

func (e *EzError) StackTrace() []StackFrame {
	frames := make([]StackFrame, 0, len(e.lineErr))
	for _, info := range e.lineErr {
		file, line := parseFileLine(info.ErrInfo)
		frames = append(frames, StackFrame{
			File:     file,
			Line:     line,
			Function: info.FnInfo,
			Message:  info.Message,
		})
	}
	return frames
}

func NewEz(str string) *EzError {
	pc, file, line, ok := runtime.Caller(1)
	var errInfo string
	if ok {
		errInfo = fmt.Sprintf("%s:%d", file, line)
	} else {
		errInfo = "unknown err file:0"
	}
	fn := runtime.FuncForPC(pc)
	info := &EzInfo{
		ErrInfo: errInfo,
		Message: str,
		Index:   1,
		AppName: zng_app.AppName,
	}
	if fn != nil {
		info.FnInfo = fn.Name()
	}
	return &EzError{
		customStr: []string{str},
		index:     1,
		lineErr:   []*EzInfo{info},
	}
}

func WrapEz(err error, customStr ...string) *EzError {
	pc, file, line, ok := runtime.Caller(1)
	var errInfo string
	if ok {
		errInfo = fmt.Sprintf("%s:%d", file, line)
	} else {
		errInfo = "unknown err file:0"
	}
	fn := runtime.FuncForPC(pc)
	info := &EzInfo{
		ErrInfo: errInfo,
		AppName: zng_app.AppName,
	}
	if fn != nil {
		info.FnInfo = fn.Name()
	}
	if err != nil {
		eErr, isEz := err.(*EzError)
		if isEz {
			info.Index = eErr.index + 1
			eErr.lineErr = append(eErr.lineErr, info)
			eErr.customStr = append(eErr.customStr, customStr...)
			eErr.index = eErr.index + 1
			return eErr
		} else {
			info.Message = err.Error()
			info.Index = 1
		}
	}
	if len(customStr) > 0 {
		info.Message = strings.Join(customStr, ",")
	}
	return &EzError{
		customStr: customStr,
		index:     1,
		lineErr:   []*EzInfo{info},
	}
}

func parseFileLine(s string) (string, int) {
	idx := strings.LastIndex(s, ":")
	if idx < 0 {
		return s, 0
	}
	line, _ := strconv.Atoi(s[idx+1:])
	return s[:idx], line
}
