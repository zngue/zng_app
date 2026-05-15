# app/nacos 配置中心

对齐中间件设计模式，提供声明式 API、配置源抽象、NacosMiddleware 链、v1/v2 SDK 兼容。

## 目录结构

```
app/nacos/
├── nacos.go         # Source/NacosMiddleware/Registry 核心类型
├── option.go        # Option + WithXxx 函数式选项
├── source_v1.go     # nacos v1 SDK 实现 (默认)
├── source_v2.go     # nacos v2 SDK 实现 (build tag: nacos_v2)
├── register.go      # NacosClient/New/RegisterParam
└── file.go          # FileSource 本地文件配置源
```

## 类型体系

| 类型 | 定义 | 对标中间件 |
|------|------|-----------|
| `Source` | `func(group, dataId string) (string, error)` | `middleware.Handler` |
| `NacosMiddleware` | `func(content string, next func(string) (string, error)) (string, error)` | `middleware.Middleware` |
| `Registry` | `struct { source, items, middlewares }` | `middleware.Registry` |

## 快速开始

### Nacos 配置源

```go
package main

import (
    "github.com/zngue/zng_app/app/nacos"
)

func main() {
    // 1. 创建 NacosClient
    client, err := nacos.New(nacos.NewOption(
        nacos.WithHost("127.0.0.1"),
        nacos.WithNamespace("develop"),
        nacos.WithUserName("nacos"),
        nacos.WithPassword("nacos"),
    ))
    if err != nil {
        panic(err)
    }

    // 2. 加载配置并反序列化
    var cfg conf.Bootstrap
    err = nacos.NewRegistry(
        nacos.WithSource(client.Source),
        nacos.WithDataId("config.yaml"),
        nacos.WithDataConfig("wechat", "config.yaml"),
    ).Scan(&cfg)
    if err != nil {
        panic(err)
    }

    // 3. 服务注册（可选）
    err = client.RegisterFunc(&nacos.RegisterParam{
        ServiceName: "go-pay-v3",
        Port:        8080,
        ClusterName: "DEFAULT",
        GroupName:   "go-pay-v3",
    })
    if err != nil {
        panic(err)
    }
}
```

### 本地文件配置源（开发环境）

```go
fileSrc := nacos.NewFileSource("./configs")

var cfg conf.Bootstrap
err := nacos.NewRegistry(
    nacos.WithSource(fileSrc.Source()),
    nacos.WithDataId("config.yaml"),
    nacos.WithDataConfig("wechat", "config.yaml"),
).Scan(&cfg)
```

目录结构：

```
configs/
├── config.yaml           # 默认 Group
└── wechat/
    └── config.yaml       # group=wechat, dataId=config.yaml
```

## API 参考

### RegistryOption

| 函数 | 说明 |
|------|------|
| `WithSource(source Source)` | 设置配置源（必须） |
| `WithDataId(dataId string)` | 添加默认 Group 的配置项 |
| `WithDataConfig(group, dataId string)` | 添加指定 Group 的配置项 |
| `WithNacosMiddleware(mw ...NacosMiddleware)` | 注册配置中间件 |

### OptionFunc

| 函数 | 默认值 | 说明 |
|------|--------|------|
| `WithNamespace(id)` | `"develop"` | 命名空间 |
| `WithHost(host)` | - | Nacos 地址 |
| `WithPort(port)` | `8848` | HTTP 端口 |
| `WithGrpcPort(port)` | `9848` | gRPC 端口（v2） |
| `WithTimeoutMs(ms)` | `5000` | 超时时间 |
| `WithUserName(name)` | - | 用户名 |
| `WithPassword(pwd)` | - | 密码 |
| `WithLogDir(dir)` | `"nacos/log"` | 日志目录 |
| `WithCacheDir(dir)` | `"nacos/cache"` | 缓存目录 |
| `WithLogLevel(level)` | - | 日志级别 |
| `WithAppendToStdout(b)` | `false` | 日志输出到 stdout |
| `WithNotLoadCacheAtStart(b)` | `false` | 启动时不加载缓存 |

### RegisterParam

| 字段 | 类型 | 说明 |
|------|------|------|
| `Port` | `int32` | 服务端口 |
| `Weight` | `float64` | 权重 |
| `ClusterName` | `string` | 集群名 |
| `ServiceName` | `string` | 服务名 |
| `GroupName` | `string` | 分组名 |

## NacosMiddleware

配置中间件，在配置加载后、反序列化前执行。可用于解密、合并、校验等。

### 执行链

```
Load → content1 + content2 + ... → NacosMiddleware 链 → YAML → map → JSON → struct
```

中间件逆序包裹，和 middleware 体系一致：

```
NacosMiddleware[0](content, next0)
  → next0 = NacosMiddleware[1](content, next1)
    → next1 = base (直接返回 content)
```

执行顺序：`NacosMiddleware[1]` → `NacosMiddleware[0]`

### 示例：解密中间件

```go
func DecryptMiddleware(content string, next func(string) (string, error)) (string, error) {
    decrypted, err := decrypt(content)
    if err != nil {
        return "", err
    }
    return next(decrypted)
}
```

### 示例：校验中间件

```go
func ValidateMiddleware(content string, next func(string) (string, error)) (string, error) {
    if content == "" {
        return "", fmt.Errorf("配置内容为空")
    }
    return next(content)
}
```

### 示例：日志中间件

```go
func LogMiddleware(content string, next func(string) (string, error)) (string, error) {
    log.Infof("加载配置内容长度: %d", len(content))
    result, err := next(content)
    if err != nil {
        log.Errorf("配置处理失败: %v", err)
    }
    return result, err
}
```

### 使用

```go
nacos.NewRegistry(
    nacos.WithSource(client.Source),
    nacos.WithDataId("config.yaml"),
    nacos.WithNacosMiddleware(
        ValidateMiddleware,
        DecryptMiddleware,
        LogMiddleware,
    ),
).Scan(&cfg)
```

## v1/v2 SDK 切换

通过 build tag 选择 Nacos SDK 版本：

```bash
# 默认使用 v1 SDK
go build ./...

# 使用 v2 SDK
go build -tags nacos_v2 ./...
```

## 自定义 Source

实现 `Source` 函数类型即可：

```go
// Etcd 配置源
type EtcdSource struct {
    Client *clientv3.Client
}

func (e *EtcdSource) Source() nacos.Source {
    return func(group string, dataId string) (string, error) {
        key := fmt.Sprintf("/%s/%s", group, dataId)
        resp, err := e.Client.Get(context.Background(), key)
        if err != nil {
            return "", err
        }
        if len(resp.Kvs) == 0 {
            return "", fmt.Errorf("配置不存在: %s", key)
        }
        return string(resp.Kvs[0].Value), nil
    }
}
```
