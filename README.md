chinadep-go-sdk
===
上海数据交易所系统接入Go SDK

## Usage
### 引入SDK
```go
import (
    "github.com/yuhu-tech/chinadep-go-sdk"
)
```
### 获取token
```go
client := chinadep.NewClient(appId, appSecret)
if err := client.ApplyToken(); err != nil {
	// ...
}
```

### 刷新token
```go
if err := client.RefreshToken(); err != nil {
	// ...
}
```

### 注册藏品
```go
client := chinadep.NewClient("appId", "appSecret")
if err := client.AssetsRegister([]*chinadep.AssetsRegisterRequest{
	{
		MetaverseAssetType: chinadep.MetaverseAssetTypeDigitalAsset,
		ChainId:            "xxx",
		// 其他参数
	},
	{
		// 可填多个资产
	},
}); err != nil {
	// ...
}

// or
assets := make([]*chinadep.AssetsRegisterRequest, 0, n)
for i := 0; i < n; i++ {
    assets = append(assets, &chinadep.AssetsRegisterRequest{
        MetaverseAssetType: chinadep.MetaverseAssetTypeDigitalAsset,
        ChainId:            "xxx",
        // 其他参数
    })
}

```