为了让大家可以快速将 [go-ioc-framework](https://github.com/jay-kim/go-ioc-framework) 项目引入并实现完整的 IOC 统一管理，我可以提供一份**标准使用指南（集成步骤）**。下面是完整详细的说明流程，包括：

------

## ✅ 第一步：引入框架

### 使用 Go Modules：

```bash
go get github.com/jay-kim/go-ioc-framework
```

或者，如果是 GitHub 镜像路径（如有同步）：

```bash
go get github.com/jay-kim/go-ioc-framework
```

------

## ✅ 第二步：初始化容器并注册服务

```go
import "github.com/jay-kim/go-ioc-framework/container"

// 初始化容器
var c = container.New()

// 注册一个构造函数（必须是无参构造函数）
c.Provide("main.UserService", func() interface{} {
    return &UserService{}
})
```

------

## ✅ 第三步：获取服务实例（自动注入）

```go
var userSvc *UserService
c.Get(&userSvc)

userSvc.Hello() // 使用服务方法
```

🔍 框架内部通过反射将实例注入到 `*userSvc` 中。

------

## ✅ 第四步（可选）：统一生命周期管理

如服务实现：

```go
type Lifecycle interface {
	Start() error
	Stop() error
}
```

容器将自动注册它，并可以统一管理：

```go
c.InitAll() // 启动所有服务
c.StopAll() // 停止所有服务
```

------

## ✅ 第五步：在项目中统一管理模块

你可以把每个模块拆成如下结构：

```go
internal/
  ├─ user/
  │   ├─ service.go      // 核心逻辑
  │   └─ api.go          // Gin 路由（或 gRPC handler）
  ├─ logger/
  │   └─ logger.go       // zap 封装
  └─ db/
      └─ mysql.go        // DB连接（带生命周期）
```

在 `main.go` 中注册所有模块：

```go
c.Provide("logger.Logger", func() interface{} {
	return logger.NewLogger()
})
c.Provide("user.UserService", func() interface{} {
	return &user.UserService{}
})
```

然后统一启动：

```go
c.InitAll()

defer c.StopAll()
```

------

## ✅ 第六步：适配 Gin / gRPC 框架

- Gin 路由挂载时：

```go
var userApi *UserApi
c.Get(&userApi)
userApi.RegisterRoutes(router)
```

- gRPC 注册服务时：

```go
var grpcSvc *GrpcUserService
c.Get(&grpcSvc)
grpcSvc.Register(grpcServer)
```

------

## 📌 总结：完整 IOC 统一管理的核心理念

| 步骤         | 说明                                            |
| ------------ | ----------------------------------------------- |
| 容器初始化   | 通过 `container.New()` 创建                     |
| 模块注册     | `Provide(name, constructor)` 注册所有模块       |
| 实例注入     | `Get(&ptr)` 获取实例并注入                      |
| 生命周期管理 | 实现 `Lifecycle` 接口统一管理 Start/Stop        |
| 框架适配     | 兼容 Gin / gRPC / zap / db 等模块统一注册与管理 |
