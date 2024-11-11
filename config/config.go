package config

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/zngue/zng_app/config/nacos"
	"gopkg.in/yaml.v3"
	"strings"
)

func NewConfig(option *nacos.CenterOptions, opts []*Opt) *Config {
	return &Config{
		Client: option.Client,
		Opts:   opts,
	}
}

type Config struct {
	Client    config_client.IConfigClient
	Opts      []*Opt
	configStr string
}

// Load
func (c *Config) Load() (err error) {
	var contentItem []string
	for _, opt := range c.Opts {
		var content string
		content, err = c.Client.GetConfig(vo.ConfigParam{
			DataId: opt.DataId,
			Group:  opt.GroupName,
		})
		if err != nil {
			return
		}
		contentItem = append(contentItem, content)
	}
	c.configStr = strings.Join(contentItem, "\n")
	return
}
func (c *Config) Scan(v any) (err error) {
	var (
		configMap = make(map[string]any)
		data      []byte
	)
	err = c.Load()
	if err != nil {
		return
	}
	if c.configStr == "" {
		err = fmt.Errorf("配置文件为空")
		return
	}
	err = yaml.Unmarshal([]byte(c.configStr), &configMap)
	if err != nil {
		return
	}
	data, err = json.Marshal(configMap)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, v)
	return
}

type Fn func(opt *Opt) *Opt

func WithDataConfig(groupName string, DataId string) Fn {
	return func(opt *Opt) *Opt {
		opt.GroupName = groupName
		opt.DataId = DataId
		return opt
	}
}
func WithDataId(dataId string) Fn {
	return func(opt *Opt) *Opt {
		opt.DataId = dataId
		return opt
	}
}

type Opt struct {
	GroupName string
	DataId    string
}

type OptionConfig struct {
	Group string
	Fns   []Fn
}

func NewOption(option *OptionConfig) (opts []*Opt) {
	for _, fn := range option.Fns {
		var o = &Opt{
			GroupName: option.Group,
		}
		fn(o)
		opts = append(opts, o)
	}
	return
}
