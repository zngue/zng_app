package nacos

import (
	"fmt"
	"log"
	"net"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// rpc 连接nacos
func NewCenterOptions() (err error) {

	sc := []constant.ServerConfig{
		{
			IpAddr: "101.43.92.81", // HTTP 可直连 IP 或域名
			Port:   8848,           // HTTP 端口
			Scheme: "http",
		},
	}
	cc := &constant.ClientConfig{
		NamespaceId:         "develop",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogLevel:            "debug",
		// EnableGrpc:          false, // ⚠️ 禁用 gRPC
	}
	cli, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  cc,
		ServerConfigs: sc,
	})
	if err != nil {
		return
	}

	var flag bool
	flag, err = cli.RegisterInstance(vo.RegisterInstanceParam{
		ServiceName: "user-service",
		Ip:          getHostIp(),
		Port:        9001,
		Ephemeral:   true,
		GroupName:   "test_user",
		Enable:      true,
		Healthy:     true,
		Weight:      10,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(flag)
	// 查询服务
	instances, err := cli.SelectAllInstances(vo.SelectAllInstancesParam{
		ServiceName: "user-service",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("发现服务实例:", instances)

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
