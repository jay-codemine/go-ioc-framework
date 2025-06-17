package grpc

import (
	"gitee.com/jay-kim/go-ioc-framework/pkg/ioc/container"
	"google.golang.org/grpc"
)

var services []func(*grpc.Server, *container.Container)

// Register 添加服务注册函数（在模块 init 中注册）
func Register(f func(*grpc.Server, *container.Container)) {
	services = append(services, f)
}

// InitServer 初始化 gRPC 服务，自动执行服务注册函数
func InitServer(c *container.Container) *grpc.Server {
	server := grpc.NewServer()

	for _, f := range services {
		f(server, c) // 调用每个模块传入的服务注册逻辑
	}

	return server
}
