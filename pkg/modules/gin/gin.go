package gin

import (
	"gitee.com/jay-kim/go-ioc-framework/pkg/ioc/container"
	"github.com/gin-gonic/gin"
)

var handlers []func(*gin.Engine, *container.Container)

// Register 添加 Gin 路由注册函数（将来统一执行）
func Register(f func(*gin.Engine, *container.Container)) {
	handlers = append(handlers, f)
}

// ✅ 新增：支持按 Group 注册路由（推荐）
func RegisterGroup(prefix string, f func(*gin.RouterGroup, *container.Container)) {
	Register(func(engine *gin.Engine, c *container.Container) {
		group := engine.Group(prefix)
		f(group, c)
	})
}

// InitServer 初始化 Gin 引擎并挂载所有注册函数
func InitServer(c *container.Container) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery())

	// 注入 ScopedContainer
	engine.Use(func(ctx *gin.Context) {
		scoped := container.NewScopedContainer(c)
		scoped.Set("gin.Context", ctx)
		ctx.Set("ioc", scoped)
		ctx.Next()
	})

	// 执行注册的 handler 函数
	for _, f := range handlers {
		f(engine, c)
	}

	return engine
}
