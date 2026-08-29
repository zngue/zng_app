package log

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	zng_app "github.com/zngue/zng_app"
)

func logF(s string, i ...any) []zap.Field {
	begin := time.Now()
	return []zap.Field{
		zap.String("serviceName", zng_app.AppName),
		zap.String("elapsed", begin.Format("2006-01-02 15:04:05")),
		zap.Any("message", fmt.Sprintf(s, i...)),
	}
}

func log(i ...any) []zap.Field {
	begin := time.Now()
	if len(i) == 0 {
		return nil
	}
	fields := []zap.Field{
		zap.String("serviceName", zng_app.AppName),
		zap.String("elapsed", begin.Format("2006-01-02 15:04:05")),
	}
	if len(i) >= 1 {
		fields = append(fields, zap.Any("message", i[0]))
	}
	return fields
}

func Errorf(s string, i ...any) {
	Default().Error("message", logF(s, i...)...)
}

func Error(i ...any) {
	Default().Error("message", log(i...)...)
}

func Warnf(s string, i ...any) {
	Default().Warn("message", logF(s, i...)...)
}

func Warn(i ...any) {
	Default().Warn("message", log(i...)...)
}

func Infof(s string, i ...any) {
	Default().Info("message", logF(s, i...)...)
}

func Info(i ...any) {
	Default().Info("message", log(i...)...)
}

func Debugf(s string, i ...any) {
	Default().Debug("message", logF(s, i...)...)
}

func Debug(i ...any) {
	Default().Debug("message", log(i...)...)
}
