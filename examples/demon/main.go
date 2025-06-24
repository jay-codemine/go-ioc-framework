package main

import (
	_ "gitee.com/jay-kim/go-ioc-framework/examples/demon/user"
	"gitee.com/jay-kim/go-ioc-framework/pkg/ioc"
	"go.uber.org/zap"
)

func main() {
	// 注册全局依赖：日志
	ioc.Provide("*zap.Logger", func() *zap.Logger {
		logger, _ := zap.NewDevelopment()
		return logger
	})

	ioc.InitAll()
	ioc.InitGinServer().Run(":8080")
}
