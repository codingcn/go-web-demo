package bootstrap

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"os"
	"gy-api/modules/tconfig"
	"gy-api/modules/zlog"
)

func Start(r *gin.Engine) {
	l := tconfig.C.GetString(fmt.Sprintf("%s.http.listen", tconfig.AppName))
	zlog.Logger.Info("服务启动成功", zap.Int("pid", os.Getpid()), zap.String("addr", "http://"+l))
	r.Run(l)
}
