package mysql

import (
	"github.com/zngue/zng_app/log"
)

type Option struct {
	Username     string
	Password     string
	Host         string
	Port         int
	Database     string
	Logger       bool
	TablePrefix  string
	LoggerLevel  int
	LoggerConfig *log.Config
}
type Fn func(opt *Option)

func DataWithUserName(username string) Fn {
	return func(opt *Option) {
		opt.Username = username
	}
}
func DataWithPassword(password string) Fn {
	return func(opt *Option) {
		opt.Password = password
	}
}
func DataWithHost(host string) Fn {
	return func(opt *Option) {
		opt.Host = host
	}
}
func DataWithPort(port int) Fn {
	return func(opt *Option) {
		opt.Port = port
	}
}
func DataWithDatabase(database string) Fn {
	return func(opt *Option) {
		opt.Database = database
	}
}
func DataWithLogger(flag bool) Fn {
	return func(opt *Option) {
		opt.Logger = flag
	}
}
func DataWithTablePrefix(prefix string) Fn {
	return func(opt *Option) {
		opt.TablePrefix = prefix
	}
}
func DataWithLoggerLevel(level int) Fn {
	return func(opt *Option) {
		opt.LoggerLevel = level
	}
}
func DataWithLoggerConfig(cfg *log.Config) Fn {
	return func(opt *Option) {
		opt.LoggerConfig = cfg
	}
}
