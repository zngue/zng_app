package log

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
)

func NewLog() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "nacos/foo.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})
	log.Println("hello")
}
