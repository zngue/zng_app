package nacos

import (
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

type NacosClient struct {
	Source       Source
	RegisterFunc RegisterFunc
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
	return &NacosClient{
		Source:       src,
		RegisterFunc: registerFn,
	}, nil
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
