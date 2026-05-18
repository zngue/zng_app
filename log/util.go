package log

import (
	"fmt"
	"time"

	"github.com/zngue/zng_app"
	"go.uber.org/zap"
)

func logF(s string, i ...any) (data []zap.Field) {
	begin := time.Now()
	data = append(data, zap.String("serviceName", zng_app.AppName))
	data = append(data, zap.String("elapsed", begin.Format("2006-01-02 15:04:05")))
	data = append(data, zap.Any("message", fmt.Sprintf(s, i...)))
	return
}
func log(i ...any) (data []zap.Field) {
	begin := time.Now()
	if len(i) == 0 {
		return
	}
	data = append(data, zap.String("serviceName", zng_app.AppName))
	data = append(data, zap.String("elapsed", begin.Format("2006-01-02 15:04:05")))
	if len(i) >= 1 {
		data = append(data, zap.Any("message", i[0]))
	}
	return
}
func Errorf(s string, i ...any) {
	data := logF(s, i...)
	Default().Error("message", data...)
}
func Error(i ...any) {
	data := log(i...)
	Default().Error("message", data...)
}
func Warnf(s string, i ...any) {
	data := logF(s, i...)
	Default().Warn("message", data...)
}
func Warn(i ...any) {
	data := log(i...)
	Default().Warn("message", data...)
}
func Infof(s string, i ...any) {
	data := logF(s, i...)
	Default().Info("message", data...)

}
func Info(i ...any) {
	data := log(i...)
	Default().Info("message", data...)
}
func Debugf(s string, i ...any) {
	data := logF(s, i...)
	Default().Debug("message", data...)
}
func Debug(i ...any) {
	data := log(i...)
	Default().Debug("message", data...)
}
