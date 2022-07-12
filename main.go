package main

//go:generate protoc --go_out=plugins=grpc:./ protos/rpc.server.proto
import (
	"github.com/gin-gonic/gin"
	"go-web-demo/app/jobs"
	"go-web-demo/app/middlewares"
	"go-web-demo/kernel/bootstrap"
	"go-web-demo/kernel/tconfig"
	"go-web-demo/kernel/zlog"
	"go-web-demo/routes"
)

func main() {
	tconfig.AppName = "demo"
	tconfig.Init()
	zlog.Init()

	// redis初始化
	//goredis.Init()

	// mysql初始化
	//torm.Init()

	go jobs.Start()
	bootstrap.Start(initRouter())
}

func initRouter() *gin.Engine {
	if tconfig.C.GetString("env") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.SetMode(gin.ReleaseMode)
	// 初始化引擎
	r := gin.New()
	// 公共中间件
	r.Use(middlewares.CommonMiddleware())
	r.Use(middlewares.ResponseLogMiddleware)

	r.Use(middlewares.CORSMiddleware())
	routes.Load(r)
	return r
}
