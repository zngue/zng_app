package log

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"io"
	"os"
	"time"
)

var DefaultLogger *zap.Logger

type SaveLog struct {
}

func (s *SaveLog) Write(p []byte) (n int, err error) {
	fmt.Println("SaveLog--Write", string(p))
	return
}

func NewLogSave() *SaveLog {
	return new(SaveLog)
}

type LevelType int8

const (
	LevelSilent LevelType = iota + 1
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
)

type Config struct {
	Filename    string
	MaxSize     int
	MaxBackups  int
	MaxAge      int
	Compress    bool
	Level       LevelType
	WriteSyncer io.Writer
}

var WriterConfigDefault = &Config{
	Filename:   "nacos/log.log",
	MaxSize:    100,
	MaxBackups: 3,
	MaxAge:     30,
	Compress:   true,
	Level:      LevelDebug,
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
		Filename:   WriterConfigDefault.Filename,   // 日志文件路径
		MaxSize:    WriterConfigDefault.MaxSize,    // 单个文件最大大小，单位为MB
		MaxBackups: WriterConfigDefault.MaxBackups, // 保留的旧文件的最大个数
		MaxAge:     WriterConfigDefault.MaxAge,     // 最大天数
		Compress:   WriterConfigDefault.Compress,   // 是否压缩
	})
	var wrSlice []zapcore.WriteSyncer
	if WriterConfigDefault.WriteSyncer != nil {
		wrSlice = append(wrSlice, zapcore.AddSync(WriterConfigDefault.WriteSyncer))
	}
	if WriterConfigDefault.Level == LevelSilent || WriterConfigDefault.Level == LevelDebug {
		wrSlice = append(wrSlice, zapcore.AddSync(os.Stdout))
	} else {
		wrSlice = append(wrSlice, fileLog)
	}
	writeSyncer := zapcore.NewMultiWriteSyncer(wrSlice...)
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

func (l *Log) Info(_ context.Context, s string, i ...any) {
	if l.LogLevel >= logger.Info {
		Default().Sugar().Info(i...)
	}
}

func (l *Log) Warn(_ context.Context, s string, i ...any) {
	if l.LogLevel >= logger.Warn {
		Default().Sugar().Warn(i...)
	}
}

func (l *Log) Error(_ context.Context, s string, i ...any) {
	if l.LogLevel >= logger.Error {
		Default().Sugar().Error(i...)
	}
}

func (l *Log) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
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

func NewLog(opt *Config) logger.Interface {
	var l = new(Log)
	if opt != nil {
		WriterConfigDefault = opt
	}
	l.LogLevel = logger.LogLevel(WriterConfigDefault.Level)
	return l
}
