以下是为你当前 `examples/basic` 示例项目量身定制的 `README.md` 文档，适合直接放在 `examples/basic/README.md` 文件中 👇：

------

```markdown
# 🧪 示例项目：go-ioc-framework 使用演示

本示例展示如何使用 [go-ioc-framework](https://gitee.com/jay-kim/go-ioc-framework) 实现服务注册、自动注入、路由绑定以及模块化管理。

------

📁 项目结构

examples/
├── basic/
│   ├── main.go              # 程序入口
│   └── README.md            # 本文件
└── internal/
└── user/
├── init.go          # 注册构造函数到容器
├── routes.go        # 路由注册
└── controller.go    # 控制器实现
```
---

## 🚀 快速开始

### 1️⃣ 进入示例目录

cd examples/basic

### 2️⃣ 启动程序

```go
go run main.go
```

控制台输出：

```go
✅ 服务启动成功，监听地址：http://localhost:8080
```

------

## 📦 访问接口

调用 `/hello` 路由查看效果：

```go
curl http://localhost:8080/user/ping
```

输出：

```json
{"message": "Hello from IOC!"}
```

------

## 🧩 模块说明

### ✅ user/controller.go

```go
type Controller struct{}

func (c *Controller) Hello(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "Hello from IOC!"})
}
```

------

### ✅ user/init.go

注册控制器构造函数到 IOC 容器：

```go
func Register(container *ioc.Container) {
	container.Provide("user.Controller", func() interface{} {
		return &Controller{}
	})
}
```

------

### ✅ user/routes.go

通过容器自动注入 controller：

```go
func RegisterRoutes(router *gin.Engine, container *ioc.Container) {
	var ctrl *Controller
	container.Get(&ctrl)

	router.GET("/hello", ctrl.Hello)
}
```

------

### ✅ main.go

```go
func main() {
	c := container.New()

	// 注册模块
	user.Register(c)

	// 初始化 HTTP 服务
	r := gin.Default()
	user.RegisterRoutes(r, c)

	// 启动容器生命周期（可选）
	c.InitAll()
	defer c.StopAll()

	// 启动服务
	r.Run(":8080")
}
```

------

## 📌 总结

- ✅ `Provide()` 注册服务构造函数
- ✅ `Get(&ptr)` 自动注入实例
- ✅ 支持 Gin / gRPC 路由绑定
- ✅ 统一生命周期管理（InitAll / StopAll）
- ✅ 清晰分层，方便拓展更多模块