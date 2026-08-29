package log

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	zng_app "github.com/zngue/zng_app"
)

type GormLog struct {
	LogLevel logger.LogLevel
}

func (l *GormLog) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *GormLog) Info(_ context.Context, s string, i ...any) {
	if l.LogLevel >= logger.Info {
		Default().Sugar().Info(i...)
	}
}

func (l *GormLog) Warn(_ context.Context, s string, i ...any) {
	if l.LogLevel >= logger.Warn {
		Default().Sugar().Warn(i...)
	}
}

func (l *GormLog) Error(_ context.Context, s string, i ...any) {
	if l.LogLevel >= logger.Error {
		Default().Sugar().Error(i...)
	}
}

func (l *GormLog) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	var data []zap.Field
	sql, rows := fc()
	data = append(data, zap.String("serviceName", zng_app.AppName))
	data = append(
		data,
		zap.String("sql", sql),
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
	)
	switch {
	case err != nil && l.LogLevel >= logger.Error:
		data = append(data, zap.Error(err))
		data = append(data, zap.String("file", utils.FileWithLineNum()))
		Default().Error("sql_info", data...)
	case l.LogLevel >= logger.Warn:
		Default().Warn("sql_info", data...)
	case l.LogLevel >= logger.Info:
		Default().Info("sql_info", data...)
	default:
		Default().Debug("sql_info", data...)
	}
}

func NewGormLog(opt *Config) logger.Interface {
	l := new(GormLog)
	if opt != nil {
		WriterConfigDefault = opt
	}
	l.LogLevel = logger.LogLevel(WriterConfigDefault.Level)
	return l
}
