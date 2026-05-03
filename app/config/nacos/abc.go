package nacos

import (
	"errors"
	"strings"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func NewNamingClient(addr, namespace string) (naming_client.INamingClient, error) {
	serverConfigs := []constant.ServerConfig{
		{IpAddr: addr, Port: 8848, GrpcPort: 49848},
	}

	clientConfig := constant.NewClientConfig(
		constant.WithNamespaceId(namespace),
		constant.WithTimeoutMs(5000),
		constant.WithLogLevel("info"),
	)

	return clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  clientConfig,
		ServerConfigs: serverConfigs,
	})
}

func waitReady(cli naming_client.INamingClient, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		// 触发内部连接逻辑
		_, err := cli.SelectInstances(vo.SelectInstancesParam{
			ServiceName: "", // 空服务名不会真正出错
			HealthyOnly: true,
		})
		if err == nil {
			// client 已经成功连接 Nacos
			return nil
		}

		// 等待期间若仍是 STARTING
		if strings.Contains(err.Error(), "STARTING") {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// 其他错误立即退出
		return err
	}
	return errors.New("nacos client not ready after timeout")
}
