package user

import (
	"gitee.com/jay-kim/go-ioc-framework/pkg/ioc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func PingHandler(ctx *gin.Context, c *ioc.Container) {
	var logger *zap.Logger
	c.Get(&logger)

	logger.Info("received ping")
	ctx.JSON(200, gin.H{"msg": "pong"})
}
