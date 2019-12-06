package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-web-demo/app/controllers/user"
	"go-web-demo/app/middlewares"
	"go-web-demo/kernel/zlog"
)

func Load(r *gin.Engine) {

	// 资源路径
	r.Static("resources/assets", "./resources/assets")

	userController := new(user.UserController)

	// 无权限路由组
	noAuthRouter := r.Group("/").Use(middlewares.NoAuth())
	{

		noAuthRouter.Any("/", func(ctx *gin.Context) {
			fmt.Println(666666)
			zlog.Logger.WithContext(ctx).Warn("测试1")
			zlog.Logger.WithContext(ctx).Error("测试2")
			ctx.String(200, "hello.")
		})

		noAuthRouter.Any("/health", func(ctx *gin.Context) {
			ctx.String(200, "ok")
		})
	}

	// 权限路由组
	authRouter := r.Group("/").Use(middlewares.JWTAuth())
	{
		// 发送短信验证码
		authRouter.POST("/user/info/bind_phone", userController.BindPhone)
	}
}
