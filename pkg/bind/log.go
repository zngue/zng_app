package bind

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/zngue/zng_app"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var DefaultLogger *zap.Logger
var LogType = "zap"

func Default() *zap.Logger {
	if DefaultLogger != nil {
		defer func() {
			defer func(logger *zap.Logger) {
				err := logger.Sync()
				if err != nil {
					fmt.Println(err)
				}
			}(DefaultLogger)
		}()
		return DefaultLogger
	}

	var wrSlice []zapcore.WriteSyncer
	wrSlice = append(wrSlice, zapcore.AddSync(os.Stdout))
	writeSyncer := zapcore.NewMultiWriteSyncer(wrSlice...)
	encoderConfig := zap.NewProductionEncoderConfig()
	level := zap.NewAtomicLevelAt(zap.InfoLevel)
	var core = zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writeSyncer,
		level,
	)
	l := zap.New(core, zap.AddCallerSkip(1))
	DefaultLogger = l
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Println(err)
		}
	}(l)
	return l
}
func logF(ctx context.Context, v any) (data []zap.Field) {
	begin := time.Now()
	data = append(data, zap.String("serviceName", zng_app.AppName))
	if val := FromUDIDContext(ctx); val != "" {
		data = append(data, zap.String(RequestIDKey, val))
	}
	if val := OperationByContext(ctx); val != "" {
		data = append(data, zap.String(OperationKey, val))
	}
	data = append(data, zap.String("elapsed", begin.Format("2006-01-02 15:04:05")))
	data = append(data, zap.Any("log_data", v))
	return
}
func Info(ctx context.Context, v any) {
	data := logF(ctx, v)
	Default().Info("app-data", data...)
}
