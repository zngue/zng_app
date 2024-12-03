package log

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"os"
	"time"
)

var DefaultLogger *zap.Logger
var LevelType int

type SaveLog struct {
}

func (s *SaveLog) Write(p []byte) (n int, err error) {
	fmt.Println("SaveLog--Write", string(p))
	return
}

func NewLogSave() *SaveLog {
	return new(SaveLog)
}

type Config struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	Level      int8
}

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
	fileLog := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "nacos/log.log", // 日志文件路径
		MaxSize:    100,             // 单个文件最大大小，单位为MB
		MaxBackups: 3,               // 保留的旧文件的最大个数
		MaxAge:     30,              // 最大天数
		Compress:   true,            // 是否压缩
	})

	writeSyncer := zapcore.NewMultiWriteSyncer(
		fileLog,
		zapcore.AddSync(NewLogSave()),
	)

	if LevelType == 5 {
		writeSyncer = zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(NewLogSave()),
		)
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	level := zap.NewAtomicLevelAt(zap.InfoLevel)
	var core zapcore.Core
	core = zapcore.NewCore(
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

type Log struct {
	LogLevel logger.LogLevel
}

func (l *Log) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *Log) Info(ctx context.Context, s string, i ...any) {
	if l.LogLevel >= logger.Info {
		Default().Sugar().Info(i...)
	}
}

func (l *Log) Warn(ctx context.Context, s string, i ...any) {
	if l.LogLevel >= logger.Warn {
		Default().Sugar().Warn(i...)
	}
}

func (l *Log) Error(ctx context.Context, s string, i ...any) {
	if l.LogLevel >= logger.Error {
		Default().Sugar().Error(i...)
	}
}

func (l *Log) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	var data []zap.Field
	sql, rows := fc()
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
		Default().Error("Error", data...)
	case l.LogLevel >= logger.Warn:
		Default().Warn("Warn", data...)
	case l.LogLevel >= logger.Info:
		Default().Info("Info", data...)
	default:
		Default().Debug("debug", data...)
	}
}

func NewLog(level int) logger.Interface {
	var l = new(Log)
	l.LogLevel = logger.LogLevel(level)
	LevelType = level
	return l
}
