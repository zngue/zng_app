package pkg

import (
	"fmt"
	"github.com/zngue/zng_app/log"
	"os"
	"strconv"
)

type DefaultConfig struct {
	Host       string
	Group      string
	Namespace  string
	HttpPort   int
	DBDatabase string
	LogLevel   string
}

type ConfigFn func(*DefaultConfig)

func DataWithHost(host string) ConfigFn {
	return func(config *DefaultConfig) {
		log.Info(fmt.Sprintf("DataWithHost,Host:%s", host))
		config.Host = host
	}
}
func DataWithGroup(group string) ConfigFn {
	return func(config *DefaultConfig) {
		log.Info(fmt.Sprintf("DataWithGroup,Group:%s", group))
		config.Group = group
	}
}
func DataWithHttpPort(port int) ConfigFn {
	return func(config *DefaultConfig) {
		log.Info(fmt.Sprintf("DataWithHttpPort,HttpPort:%d", port))
		config.HttpPort = port
	}
}
func DataWithDBDatabase(database string) ConfigFn {
	return func(config *DefaultConfig) {
		log.Info(fmt.Sprintf("DataWithDBDatabase,DBDatabase:%s", database))
		config.DBDatabase = database
	}
}
func DataWithNamespace(namespace string) ConfigFn {
	return func(config *DefaultConfig) {
		log.Info(fmt.Sprintf("DataWithNamespace,Namespace:%s", namespace))
		config.Namespace = namespace
	}
}
func DataWithLogLevel(level string) ConfigFn {
	return func(config *DefaultConfig) {
		log.Info(fmt.Sprintf("DataWithLogLevel,LogLevel:%s", level))
		config.LogLevel = level
	}
}

func NewConfig(fns ...ConfigFn) (config *DefaultConfig, err error) {
	config = &DefaultConfig{
		Host:      "nacos.zngue.com",
		Group:     "common",
		Namespace: "develop",
		HttpPort:  16667,
		LogLevel:  "info", //info debug warn error
	}
	var host = os.Getenv("HOST")
	if host != "" {
		config.Host = host
	}
	log.Info(fmt.Sprintf("_HOST,Host:%s", config.Host))
	var dbGroupName = os.Getenv("DB_GROUP")
	if dbGroupName != "" {
		config.Group = dbGroupName
	}
	log.Info(fmt.Sprintf("DB_GROUP,Group:%s", config.Group))
	var namespace = os.Getenv("NAMESPACE")
	if namespace != "" {
		config.Namespace = namespace
	}
	log.Info(fmt.Sprintf("_NAMESPACE,Namespace:%s", config.Namespace))
	var port = os.Getenv("HTTP_PORT")
	if port != "" {
		var httpPort, _ = strconv.Atoi(port)
		if httpPort > 0 {
			config.HttpPort = httpPort
		}
	}
	log.Info(fmt.Sprintf("HttpPort:%d", config.HttpPort))
	database := os.Getenv("DB_DATABASE")
	if database != "" {
		config.DBDatabase = database
	}
	log.Info(fmt.Sprintf("DB_DATABASE,DBDatabase:%s", config.DBDatabase))
	for _, fn := range fns {
		fn(config)
	}
	if config.Host == "" {
		err = fmt.Errorf("HOST 不能为空")
		return
	}
	if config.Group == "" {
		err = fmt.Errorf("DB_GROUP 不能为空")
		return
	}
	if config.Namespace == "" {
		err = fmt.Errorf("NAMESPACE 不能为空")
		return
	}
	log.Info(fmt.Sprintf("NewConfig,config:%+v", config))
	return
}
