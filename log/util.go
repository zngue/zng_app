package log

import (
	"fmt"
	"github.com/zngue/zng_app"
	"go.uber.org/zap"
	"time"
)

func logF(s string, i ...any) (data []zap.Field) {
	begin := time.Now()
	data = append(data, zap.String("serviceName", zng_app.AppName))
	data = append(data, zap.String("elapsed", begin.Format("2006-01-02 15:04:05")))
	data = append(data, zap.Any("log_data", fmt.Sprintf(s, i...)))
	return
}
func log(i ...any) (data []zap.Field) {
	begin := time.Now()
	if len(i) == 0 {
		return
	}
	data = append(data, zap.String("serviceName", zng_app.AppName))
	data = append(data, zap.String("elapsed", begin.Format("2006-01-02 15:04:05")))
	if len(i) > 0 {
		for index, val := range i {
			if index == 0 {
				data = append(data, zap.Any(fmt.Sprintf("log_data"), val))
			} else {
				data = append(data, zap.Any(fmt.Sprintf("log_data_%d", index), val))
			}
		}
	}
	return
}
func Errorf(s string, i ...any) {
	data := logF(s, i...)
	Default().Error("app_info", data...)
}
func Error(i ...any) {
	data := log(i...)
	Default().Error("app_info", data...)
}
func Warnf(s string, i ...any) {
	data := logF(s, i...)
	Default().Warn("app_info", data...)
}
func Warn(i ...any) {
	data := log(i...)
	Default().Warn("app_info", data...)
}
func Infof(s string, i ...any) {
	data := logF(s, i...)
	Default().Info("app_info", data...)

}
func Info(i ...any) {
	data := log(i...)
	Default().Info("app_info", data...)
}
func Debugf(s string, i ...any) {
	data := logF(s, i...)
	Default().Debug("app_info", data...)
}
func Debug(i ...any) {
	data := log(i...)
	Default().Debug("app_info", data...)
}
