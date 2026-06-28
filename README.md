# anox-sdk

Standalone Go SDK for Anox service registration, heartbeat, config pull, and log reporting.

## Install

```bash
go get github.com/mtit/anox-sdk@latest
```

## Usage

```go
package main

import sdk "github.com/mtit/anox-sdk"

func main() {
	client, err := sdk.NewClient(sdk.Config{
		AnoxURL:     "ws://127.0.0.1:8848/ws",
		ServiceName: "demo-service",
		HttpHost:    "127.0.0.1",
		HttpPort:    "8080",
	})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	client.Logger().LogInfo("startup", "service online", nil, true)
}
```

Before publishing, replace the module path in `go.mod` with your real GitHub repository path.
