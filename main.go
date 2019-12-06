package main

//go:generate protoc --go_out=plugins=grpc:./ protos/rpc.server.proto
import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"go-web-demo/app/jobs"
	"go-web-demo/app/middlewares"
	"go-web-demo/config"
	"go-web-demo/kernel/fixs"
	"go-web-demo/kernel/zlog"
	"go-web-demo/routes"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	//"github.com/gin-contrib/gzip"
	"github.com/cloudflare/tableflip"
	"github.com/gin-gonic/gin"
)

func main() {
	if config.C.Env != "local" {
		// 本地不重定向std
		setupStdFile()
	}

	// 数据库链接初始化
	//gorm.InitDbConnect()
	// redis链接初始化
	//goredis.InitConnect()
	// go程启动计划脚本
	go jobs.Start()
	// TODO: 由于gin官方暂时没有支持v9，所有这是v9的临时更新方案,等gin更新升级后可移除
	binding.Validator = new(fixs.DefaultValidator)

	Graceful()
}

func initRouter() *gin.Engine {
	if config.C.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.SetMode(gin.ReleaseMode)
	// 初始化引擎
	r := gin.New()
	r.Use(middlewares.Metric())
	// 公共中间件
	r.Use(middlewares.CommonMiddleware())

	r.Use(middlewares.CORSMiddleware())
	routes.Load(r)
	return r
}

func Graceful() {
	if config.C.PIDFile == "" {
		log.Fatalln("pidFile not found!")
	}
	pidFile := config.C.PIDFile
	upg, err := tableflip.New(tableflip.Options{
		PIDFile: pidFile,
	})
	if err != nil {
		panic(err)
	}
	defer upg.Stop()

	// Do an upgrade on SIGHUP
	var exit bool
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGUSR2, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
		for s := range sig {
			switch s {
			case syscall.SIGHUP, syscall.SIGUSR2:
				log.Println("Upgrade start:", s)
				err := upg.Upgrade()
				if err != nil {
					log.Println("Upgrade failed:", err)
					continue
				}
				log.Println("Upgrade succeeded")
			case syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT:
				log.Println("Shutdown Server start:", s)
				exit = true
				upg.Stop()
				log.Println("Shutdown Server ...")
			}
		}
	}()

	ln, err := upg.Fds.Listen("tcp", config.C.Http.Listen)
	if err != nil {
		log.Fatalln("Can't listen:", err)
	}

	server := &http.Server{
		Handler: initRouter(),
	}

	go func() {
		err := server.Serve(ln)
		if err != http.ErrServerClosed {
			log.Println("HTTP server:", err)
		}
	}()
	zlog.Logger.Info("服务启动成功", zap.Int("pid", os.Getpid()), zap.String("addr", "http://"+config.C.Http.Listen))

	err = ioutil.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0755)
	if err != nil {
		panic(err)
	}
	if err := upg.Ready(); err != nil {
		panic(err)
	}
	<-upg.Exit()

	// Make sure to set a deadline on exiting the process
	// after upg.Exit() is closed. No new upgrades can be
	// performed if the parent doesn't exit.
	time.AfterFunc(20*time.Second, func() {
		log.Println("Graceful shutdown timed out")
		os.Exit(1)
	})
	// Wait for connections to drain.
	server.Shutdown(context.Background())

	if exit {
		_ = os.Remove(pidFile)
	}
}

func setupStdFile() {
	outFile, err := os.OpenFile(config.C.Log.LogFileNameStdout, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatalf(fmt.Sprintf("open file error:%s|%s", config.C.Log.LogFileNameStdout, err.Error()))

	}
	errFile, err := os.OpenFile(config.C.Log.LogFileNameStderr, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatalf(fmt.Sprintf("open file error:%s|%s", config.C.Log.LogFileNameStderr, err.Error()))
	}
	os.Stdout = outFile
	os.Stderr = errFile
}
