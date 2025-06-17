package ioc

import (
	"gitee.com/jay-kim/go-ioc-framework/pkg/ioc/container"

	// 模块插件（注意：为了避免与官方包名冲突，起别名）
	ginmod "gitee.com/jay-kim/go-ioc-framework/pkg/modules/gin"
	grpcmod "gitee.com/jay-kim/go-ioc-framework/pkg/modules/grpc"

	// 官方框架包引用（建议起别名以避免冲突）
	ginpkg "github.com/gin-gonic/gin"
	gogrpc "google.golang.org/grpc"
)

// 全局容器实例（框架级）
var GlobalContainer = container.New()

// Provide 注册服务构造函数
func Provide(name string, constructor interface{}) {
	GlobalContainer.Provide(name, constructor)
}

// Get 获取已注册的服务实例
func Get(target interface{}) {
	GlobalContainer.Get(target)
}

// InitAll 初始化所有服务（执行构造函数、Start()）
func InitAll() {
	GlobalContainer.InitAll()
}

// StopAll 销毁所有服务（调用 Stop()）
func StopAll() {
	GlobalContainer.StopAll()
}

// RegisterGRPCService 注册 gRPC 服务（外部调用）
func RegisterGRPCService(fn func(server *gogrpc.Server, c *container.Container)) {
	grpcmod.Register(fn)
}

// InitGRPCServer 初始化并返回 gRPC server
func InitGRPCServer() *gogrpc.Server {
	return grpcmod.InitServer(GlobalContainer)
}

// RegisterGinGroup 提供给外部模块使用，自动封装路由分组注册
func RegisterGinGroup(prefix string, fn func(*ginpkg.RouterGroup, *container.Container)) {
	ginmod.RegisterGroup(prefix, fn)
}

// RegisterGinHandler 注册 Gin 路由处理函数（外部调用）
func RegisterGinHandler(fn func(router *ginpkg.Engine, c *container.Container)) {
	ginmod.Register(fn)
}

// InitGinServer 初始化并返回 Gin 引擎
func InitGinServer() *ginpkg.Engine {
	return ginmod.InitServer(GlobalContainer)
}

// 可选：导出容器类型（便于外部项目类型兼容）
type Container = container.Container
