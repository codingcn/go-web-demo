//go:build !plan9 && !windows
// +build !plan9,!windows

package bootstrap

import (
	"context"
	"github.com/cloudflare/tableflip"
	"github.com/gin-gonic/gin"
	"go-web-demo/kernel/tconfig"
	"go-web-demo/kernel/zlog"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func Start(r *gin.Engine) {
	if tconfig.C.GetString("pid_file") == "" {
		log.Fatalln("pidFile not found!")
	}
	pidFile := tconfig.C.GetString("pid_file")
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
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
		for s := range sig {
			switch s {
			case syscall.SIGHUP:
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

	l := tconfig.C.GetString("http.listen")
	ln, err := upg.Fds.Listen("tcp", l)
	if err != nil {
		log.Fatalln("Can't listen:", err)
	}

	server := &http.Server{
		Handler: r,
	}

	go func() {
		err := server.Serve(ln)
		if err != http.ErrServerClosed {
			log.Println("HTTP server:", err)
		}
	}()
	zlog.Logger.Info("服务启动成功", zap.Int("pid", os.Getpid()), zap.String("addr", "http://"+l))

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
