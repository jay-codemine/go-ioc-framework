package ioc

import (
	ginpkg "github.com/gin-gonic/gin"
	"github.com/jay-kim/go-ioc/pkg/ioc/container"
	"github.com/jay-kim/go-ioc/pkg/modules/gin"
	"github.com/jay-kim/go-ioc/pkg/modules/grpc"
	"google.golang.org/grpc"
)

var GlobalContainer = container.New()

// Provide 注册一个服务构造函数
func Provide(name string, constructor interface{}) {
	GlobalContainer.Provide(name, constructor)
}

// Get 获取已注册的服务实例（使用指针接收）
func Get(target interface{}) {
	GlobalContainer.Get(target)
}

// InitAll 初始化所有服务（会自动执行构造函数）
func InitAll() {
	GlobalContainer.InitAll()
}

// StopAll 优雅销毁所有服务
func StopAll() {
	GlobalContainer.StopAll()
}

// RegisterGRPCService 注册一个 gRPC 服务（由使用者提供注册逻辑）
func RegisterGRPCService(fn func(server *grpc.Server, c *container.Container)) {
	grpc.Register(fn)
}

// InitGRPCServer 初始化并返回一个 gRPC Server
func InitGRPCServer() *grpc.Server {
	return grpc.InitServer(GlobalContainer)
}

// RegisterGinHandler 注册一个 Gin 路由处理器
func RegisterGinHandler(fn func(router *ginpkg.Engine, c *container.Container)) {
	gin.Register(fn)
}

// InitGinServer 初始化并返回一个 Gin Engine
func InitGinServer() *ginpkg.Engine {
	return gin.InitServer(GlobalContainer)
}
