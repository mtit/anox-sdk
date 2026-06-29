# anox-sdk

一个独立的 Go SDK，用于 Anox 服务注册、心跳上报、配置拉取和日志上报。

## 安装

```bash
go get github.com/mtit/anox-sdk@latest
```

## 使用示例

```go
package main

import sdk "github.com/mtit/anox-sdk"

func main() {
	client, err := sdk.NewClient(sdk.Config{
		AnoxURL:     "ws://127.0.0.1:8848/ws",
		ServiceName: "demo-service",
		HttpPort:    "8080",
	})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	config := client.Config()
	port := config.GetString("server.port")
	debug := config.GetBool("debug", false)
	workerCount := config.GetInt("worker.count", 4)

	_ = port
	_ = debug
	_ = workerCount

	client.Logger().LogInfo("startup", "service online", nil, true)
}
```

## 获取配置

SDK 会在注册成功后拉取全局配置（`_global`）和当前服务配置（`ServiceName`）。`NewClient` 返回前会等待首轮配置拉取完成，后续在 Anox 推送更新时自动刷新。

```go
config := client.Config()

// 优先读取服务级配置，若不存在则回退到全局配置。
dbHost := config.GetString("db.host")

// 带默认值的类型化读取。
timeoutMs := config.GetInt("http.timeout_ms", 3000)
rateLimit := config.GetInt64("rate.limit", 1000)
sampleRatio := config.GetFloat64("trace.sample_ratio", 0.1)
debug := config.GetBool("debug", false)

// 只读取当前服务下的配置。
serviceOnly := config.GetServiceString("feature.flag")

// 只读取全局配置。
globalOnly := config.GetGlobalString("region")
```

配置读取规则：

- `GetString(key)` / `GetInt(...)` / `GetInt64(...)` / `GetFloat64(...)` / `GetBool(...)`：优先读取服务级配置，不存在时回退到全局配置。
- `GetServiceString(key)`：只读取当前服务下的配置。
- `GetGlobalString(key)`：只读取 `_global` 下的配置。
- 类型化读取方法在配置不存在或解析失败时，返回传入的默认值。
