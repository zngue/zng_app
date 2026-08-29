package nacos

import (
	"context"
	"fmt"
	"net"
	"os"
)

type RegisterParam struct {
	Port        int32
	Weight      float64
	ClusterName string
	ServiceName string
	GroupName   string
}

type RegisterFunc func(params *RegisterParam) error

type SelectInstanceFunc func(ctx context.Context, serviceName, groupName string) (string, error)

type PublishFunc func(group, dataId, content string) (bool, error)

type NacosClient struct {
	Source         Source
	RegisterFunc   RegisterFunc
	SelectInstance SelectInstanceFunc
	PublishConfig  PublishFunc
}

func New(opt *Option) (*NacosClient, error) {
	src, err := newSource(opt)
	if err != nil {
		return nil, err
	}
	registerFn, err := newRegisterFunc(opt)
	if err != nil {
		return nil, err
	}
	selectFn, err := newSelectInstance(opt)
	if err != nil {
		return nil, err
	}
	publishFn, err := newPublishFunc(opt)
	if err != nil {
		return nil, err
	}
	return &NacosClient{
		Source:         src,
		RegisterFunc:   registerFn,
		SelectInstance: selectFn,
		PublishConfig:  publishFn,
	}, nil
}

func (c *NacosClient) Resolver() SelectInstanceFunc {
	return c.SelectInstance
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

func getHostName() (string, error) {
	return os.Hostname()
}
