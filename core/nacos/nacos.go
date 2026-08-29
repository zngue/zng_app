package nacos

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

type Source func(group string, dataId string) (string, error)

type NacosMiddleware func(content string, next func(string) (string, error)) (string, error)

type LoadItem struct {
	Group  string
	DataId string
}

type Registry struct {
	source      Source
	items       []LoadItem
	middlewares []NacosMiddleware
}

type RegistryOption func(*Registry)

func WithSource(source Source) RegistryOption {
	return func(r *Registry) {
		r.source = source
	}
}

func WithDataId(dataId string) RegistryOption {
	return func(r *Registry) {
		r.items = append(r.items, LoadItem{DataId: dataId})
	}
}

func WithDataConfig(group string, dataId string) RegistryOption {
	return func(r *Registry) {
		r.items = append(r.items, LoadItem{Group: group, DataId: dataId})
	}
}

func WithNacosMiddleware(mw ...NacosMiddleware) RegistryOption {
	return func(r *Registry) {
		r.middlewares = append(r.middlewares, mw...)
	}
}

func NewRegistry(opts ...RegistryOption) *Registry {
	r := &Registry{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *Registry) Scan(v any) (err error) {
	content, err := r.Load()
	if err != nil {
		return
	}
	if content == "" {
		return fmt.Errorf("配置文件为空")
	}
	content, err = r.applyMiddlewares(content)
	if err != nil {
		return
	}
	var configMap = make(map[string]any)
	err = yaml.Unmarshal([]byte(content), &configMap)
	if err != nil {
		return
	}
	var data []byte
	data, err = json.Marshal(configMap)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, v)
	return
}

func (r *Registry) Load() (content string, err error) {
	if r.source == nil {
		return "", fmt.Errorf("配置源未设置，请使用 WithSource")
	}
	var contentItem []string
	for _, item := range r.items {
		var c string
		c, err = r.source(item.Group, item.DataId)
		if err != nil {
			return
		}
		contentItem = append(contentItem, c)
	}
	content = strings.Join(contentItem, "\n")
	return
}

func (r *Registry) applyMiddlewares(content string) (string, error) {
	h := func(s string) (string, error) { return s, nil }
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		h = wrapNacosMiddleware(r.middlewares[i], h)
	}
	return h(content)
}

func wrapNacosMiddleware(mw NacosMiddleware, next func(string) (string, error)) func(string) (string, error) {
	return func(content string) (string, error) {
		return mw(content, next)
	}
}

type LogType string

const (
	LogError LogType = "error"
	LogDebug LogType = "debug"
	LogInfo  LogType = "info"
	LogWarn  LogType = "warn"
)

func (l LogType) ZapLevel() zapcore.Level {
	switch l {
	case LogDebug:
		return zapcore.DebugLevel
	case LogInfo:
		return zapcore.InfoLevel
	case LogWarn:
		return zapcore.WarnLevel
	case LogError:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
