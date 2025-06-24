package user

import "gitee.com/jay-kim/go-ioc-framework/pkg/ioc"

func init() {
	ioc.RegisterGinGroup("/user", RegisterRoutes)
}
