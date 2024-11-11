package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"net"
)

type Fn func(*Option) *Option

type LogType string

// debug,info,warn,error
const (
	ERROR LogType = "error"
	DEBUG LogType = "debug"
	INFO  LogType = "info"
	WARN  LogType = "warn"
)

func DataWithNamespace(namespaceId string) Fn {
	return func(option *Option) *Option {
		option.NamespaceId = namespaceId
		return option
	}
}
func DataWithLogDir(s string) Fn {
	return func(option *Option) *Option {
		option.LogDir = s
		return option
	}
}
func DataWithCacheDir(s string) Fn {
	return func(option *Option) *Option {
		option.CacheDir = s
		return option
	}
}
func DataWithLogLevel(s LogType) Fn {
	return func(option *Option) *Option {
		option.LogLevel = s
		return option
	}
}
func DataWithAppendToStdout(b bool) Fn {
	return func(option *Option) *Option {
		option.AppendToStdout = b
		return option
	}
}
func DataWithUserName(userName string) Fn {
	return func(option *Option) *Option {
		option.Username = userName
		return option
	}
}
func DataWithPassword(password string) Fn {
	return func(option *Option) *Option {
		option.Password = password
		return option
	}
}
func DataWithHost(host string) Fn {
	return func(option *Option) *Option {
		option.Host = host
		return option
	}
}
func DataWithPort(port int) Fn {
	return func(option *Option) *Option {
		option.Port = port
		return option
	}
}
func DataWithTimeoutMs(timeoutMs uint64) Fn {
	return func(option *Option) *Option {
		option.TimeoutMs = timeoutMs
		return option
	}
}

func NewOption(fns ...Fn) (opt *Option) {
	opt = &Option{
		NamespaceId:    "develop",
		LogDir:         "nacos/log",
		CacheDir:       "nacos/cache",
		TimeoutMs:      5000,
		AppendToStdout: false,
		Port:           8848,
	}
	for _, fn := range fns {
		fn(opt)
	}
	return
}

type Option struct {
	NamespaceId         string
	LogDir              string
	CacheDir            string
	LogLevel            LogType
	AppendToStdout      bool
	Username            string
	Password            string
	Host                string
	Port                int
	TimeoutMs           uint64
	NotLoadCacheAtStart bool
}
type RegisterInstanceParam struct {
	Port        int32   //required
	Weight      float64 //required,it must be lager than 0
	ClusterName string  //optional,default:DEFAULT
	ServiceName string  //required
	GroupName   string  //optional,default:DEFAULT_GROUP
}

type OptionLoad struct {
}

type CenterOptions struct {
	opt    *Option
	params vo.NacosClientParam
	Client config_client.IConfigClient
	Server naming_client.INamingClient
}

func (c *CenterOptions) clientConfig() (err error) {
	c.Client, err = clients.NewConfigClient(c.params)
	return
}
func (c *CenterOptions) clientServer() (err error) {
	c.Server, err = clients.NewNamingClient(c.params)
	return
}
func (c *CenterOptions) RegisterInstance(params *RegisterInstanceParam) (err error) {
	var flag bool
	flag, err = c.Server.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          getHostIp(), //使用本机ip
		Port:        uint64(params.Port),
		ServiceName: params.ServiceName,
		GroupName:   params.GroupName,
		ClusterName: params.ClusterName,
		Weight:      params.Weight,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})
	if err != nil {
		return
	}
	if !flag {
		err = fmt.Errorf("服务注册失败 ServiceName:%s,GroupName：%s", params.ServiceName, params.GroupName)
	}
	return
}

func NewCenterOptions(opt *Option) (options *CenterOptions, err error) {
	cc := &constant.ClientConfig{
		NamespaceId:         opt.NamespaceId,
		TimeoutMs:           opt.TimeoutMs,
		NotLoadCacheAtStart: opt.NotLoadCacheAtStart,
		LogDir:              opt.LogDir,
		CacheDir:            opt.CacheDir,
		LogLevel:            string(opt.LogLevel),
		AppendToStdout:      opt.AppendToStdout,
		Username:            opt.Username,
		Password:            opt.Password,
	}
	sc := constant.ServerConfig{
		IpAddr: opt.Host,
		Port:   uint64(opt.Port),
	}
	options = &CenterOptions{
		opt: opt,
		params: vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: []constant.ServerConfig{sc},
		},
	}
	err = options.clientConfig()
	if err != nil {
		return
	}
	err = options.clientServer()
	if err != nil {
		return
	}
	return
}
func getHostIp() string {
	addrList, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("get current host ip err: ", err)
		return ""
	}
	var ip string
	for _, address := range addrList {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
				break
			}
		}
	}
	return ip
}
