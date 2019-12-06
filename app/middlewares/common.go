package middlewares

import (
	"bytes"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	uuid "github.com/satori/go.uuid"
	"go-web-demo/app/helpers"
	"go-web-demo/kernel/zlog"
	"go.uber.org/zap"
	"io/ioutil"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func CommonMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 每次请求都初始化一次配置
		u1 := uuid.NewV4()
		trace := helpers.Md5(u1.String())
		zlog.Logger.NewContext(ctx, zap.String("trace", trace))
		zlog.Logger.NewContext(ctx, zap.String("request_method", ctx.Request.Method))
		headers, _ := json.Marshal(ctx.Request.Header)
		zlog.Logger.NewContext(ctx, zap.String("request_headers", string(headers)))
		zlog.Logger.NewContext(ctx, zap.String("request_url", ctx.Request.URL.String()))
		bodyBytes, _ := ctx.GetRawData()
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // 关键点
		zlog.Logger.NewContext(ctx, zap.String("request_params", string(bodyBytes)))
		ctx.Next()
	}
}
