//go:build !nacos_v2

package nacos

import (
	"context"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func newSource(opt *Option) (Source, error) {
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
	client, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  cc,
		ServerConfigs: []constant.ServerConfig{sc},
	})
	if err != nil {
		return nil, err
	}
	return func(group string, dataId string) (string, error) {
		return client.GetConfig(vo.ConfigParam{DataId: dataId, Group: group})
	}, nil
}

func newPublishFunc(opt *Option) (PublishFunc, error) {
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
	client, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  cc,
		ServerConfigs: []constant.ServerConfig{sc},
	})
	if err != nil {
		return nil, err
	}
	return func(group, dataId, content string) (bool, error) {
		return client.PublishConfig(vo.ConfigParam{
			DataId:  dataId,
			Group:   group,
			Content: content,
		})
	}, nil
}

func newRegisterFunc(opt *Option) (RegisterFunc, error) {
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
	client, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  cc,
		ServerConfigs: []constant.ServerConfig{sc},
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

func newSelectInstance(opt *Option) (SelectInstanceFunc, error) {
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
	client, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  cc,
		ServerConfigs: []constant.ServerConfig{sc},
	})
	if err != nil {
		return nil, err
	}
	return func(ctx context.Context, serviceName, groupName string) (string, error) {
		instances, err := client.SelectInstances(vo.SelectInstancesParam{
			ServiceName: serviceName,
			GroupName:   groupName,
			HealthyOnly: true,
		})
		if err != nil {
			return "", err
		}
		if len(instances) == 0 {
			return "", fmt.Errorf("no healthy instance for %s/%s", groupName, serviceName)
		}
		inst := instances[0]
		return fmt.Sprintf("%s:%d", inst.Ip, inst.Port), nil
	}, nil
}
