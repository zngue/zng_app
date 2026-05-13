package errors_ez

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	"github.com/zngue/zng_app"
)

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

var defaultError = "unknown err file:0"

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
func New(str string) *EzError {
	pc, file, line, ok := runtime.Caller(1)
	var errInfo string
	if ok {
		errInfo = fmt.Sprintf("%s:%d", file, line)
	} else {
		errInfo = defaultError
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

// Wrap wraps an error with caller info. Note: when err is already an *EzError,
// it modifies the original *EzError in place and returns it.
func Wrap(err error, customStr ...string) *EzError {
	pc, file, line, ok := runtime.Caller(1)
	var errInfo string
	if ok {
		errInfo = fmt.Sprintf("%s:%d", file, line)
	} else {
		errInfo = defaultError
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
