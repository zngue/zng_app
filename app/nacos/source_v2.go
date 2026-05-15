//go:build nacos_v2

package nacos

import (
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func newSource(opt *Option) (Source, error) {
	cc := constant.NewClientConfig(
		constant.WithNamespaceId(opt.NamespaceId),
		constant.WithTimeoutMs(opt.TimeoutMs),
		constant.WithNotLoadCacheAtStart(opt.NotLoadCacheAtStart),
		constant.WithLogDir(opt.LogDir),
		constant.WithCacheDir(opt.CacheDir),
		constant.WithLogLevel(string(opt.LogLevel)),
		constant.WithUsername(opt.Username),
		constant.WithPassword(opt.Password),
	)
	sc := []constant.ServerConfig{
		{
			IpAddr:   opt.Host,
			Port:     uint64(opt.Port),
			GrpcPort: uint64(opt.GrpcPort),
		},
	}
	client, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  cc,
		ServerConfigs: sc,
	})
	if err != nil {
		return nil, err
	}
	return func(group string, dataId string) (string, error) {
		return client.GetConfig(vo.ConfigParam{DataId: dataId, Group: group})
	}, nil
}

func newRegisterFunc(opt *Option) (RegisterFunc, error) {
	cc := constant.NewClientConfig(
		constant.WithNamespaceId(opt.NamespaceId),
		constant.WithTimeoutMs(opt.TimeoutMs),
		constant.WithNotLoadCacheAtStart(opt.NotLoadCacheAtStart),
		constant.WithLogDir(opt.LogDir),
		constant.WithCacheDir(opt.CacheDir),
		constant.WithLogLevel(string(opt.LogLevel)),
		constant.WithUsername(opt.Username),
		constant.WithPassword(opt.Password),
	)
	sc := []constant.ServerConfig{
		{
			IpAddr:   opt.Host,
			Port:     uint64(opt.Port),
			GrpcPort: uint64(opt.GrpcPort),
		},
	}
	client, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  cc,
		ServerConfigs: sc,
	})
	if err != nil {
		return nil, err
	}
	return func(params *RegisterParam) error {
		ip := getHostIp()
		name, _ := getHostName()
		flag, err := client.RegisterInstance(vo.RegisterInstanceParam{
			Ip:          ip,
			Port:        uint64(params.Port),
			ServiceName: params.ServiceName,
			GroupName:   params.GroupName,
			ClusterName: params.ClusterName,
			Weight:      params.Weight,
			Enable:      true,
			Healthy:     true,
			Ephemeral:   true,
			Metadata: map[string]string{
				"port":     fmt.Sprintf("%d", params.Port),
				"hostName": name,
			},
		})
		if err != nil {
			return err
		}
		if !flag {
			return fmt.Errorf("服务注册失败 ServiceName:%s GroupName:%s", params.ServiceName, params.GroupName)
		}
		return nil
	}, nil
}
