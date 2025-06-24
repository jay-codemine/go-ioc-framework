package user

import (
	"gitee.com/jay-kim/go-ioc-framework/pkg/ioc"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, c *ioc.Container) {
	rg.GET("/ping", func(ctx *gin.Context) {
		PingHandler(ctx, c)
	})
}
