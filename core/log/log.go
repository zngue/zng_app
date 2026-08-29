package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	zng_app "github.com/zngue/zng_app"
)

var DefaultLogger *zap.Logger

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
	ProjectName string
}

var WriterConfigDefault = &Config{
	Filename:    "/project.log",
	ProjectName: zng_app.AppName,
	MaxSize:     100,
	MaxBackups:  3,
	MaxAge:      30,
	Compress:    true,
	Level:       LevelDebug,
}

func WriteSyncerInfo(logger *zap.Logger) {
	if WriterConfigDefault.WriteSyncer != nil {
		err := logger.Sync()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func Default() *zap.Logger {
	if DefaultLogger != nil {
		defer WriteSyncerInfo(DefaultLogger)
		return DefaultLogger
	}
	fileLog := zapcore.AddSync(ZapLoggerWriter())
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
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writeSyncer,
		level,
	)
	DefaultLogger = zap.New(core, zap.AddCallerSkip(1))
	return DefaultLogger
}

func ZapLoggerWriter() io.Writer {
	filename := WriterConfigDefault.Filename
	hook, err := rotatelogs.New(
		strings.Replace(filename, ".log", "", -1)+"-%Y%m%d.log",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*time.Duration(WriterConfigDefault.MaxAge)),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		panic(err)
	}
	return hook
}
