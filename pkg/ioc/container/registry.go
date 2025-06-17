package container

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var (
	ginHandlers  []func(*gin.Engine, *Container)
	grpcServices []func(*grpc.Server, *Container)
)

// RegisterGinHandler 注册一个 Gin 路由挂载函数
func RegisterGinHandler(f func(*gin.Engine, *Container)) {
	ginHandlers = append(ginHandlers, f)
}

// ApplyGinHandlers 调用所有已注册的 Gin 路由挂载函数
func ApplyGinHandlers(engine *gin.Engine, c *Container) {
	for _, f := range ginHandlers {
		f(engine, c)
	}
}

// RegisterGRPCService 注册一个 gRPC 服务挂载函数
func RegisterGRPCService(f func(*grpc.Server, *Container)) {
	grpcServices = append(grpcServices, f)
}

// ApplyGRPCServices 调用所有已注册的 gRPC 服务挂载函数
func ApplyGRPCServices(server *grpc.Server, c *Container) {
	for _, f := range grpcServices {
		f(server, c)
	}
}
