package middlewares

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	uuid "github.com/satori/go.uuid"
	"go-web-demo/app/helpers"
	"go-web-demo/kernel/zlog"
	"go.uber.org/zap"
	"io/ioutil"
	"strconv"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func CommonMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 每次请求都初始化一次配置
		u1 := uuid.NewV4()
		trace := helpers.Md5(u1.String())
		zlog.Logger.NewGinContext(ctx, zap.String("trace", trace))
		zlog.Logger.NewGinContext(ctx, zap.String("request_method", ctx.Request.Method))
		headers, _ := json.Marshal(ctx.Request.Header)
		zlog.Logger.NewGinContext(ctx, zap.String("request_headers", string(headers)))
		zlog.Logger.NewGinContext(ctx, zap.String("request_url", ctx.Request.URL.String()))
		if ctx.Request.Body != nil && ctx.Request.RequestURI != "/article/upload" {
			bodyBytes, _ := ctx.GetRawData()
			ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // 关键点
			zlog.Logger.NewGinContext(ctx, zap.String("request_params", string(bodyBytes)))
		}
		ctx.Next()
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func ResponseLogMiddleware(ctx *gin.Context) {
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
	ctx.Writer = blw
	startTime := time.Now()
	ctx.Next()
	duration := int(time.Since(startTime).Milliseconds())
	ctx.Header("X-Response-Time", strconv.Itoa(duration))
	if ctx.Request.URL.Path == "/cos/object" && ctx.Request.Method == "POST" {
		return
	}
	zlog.Logger.WithGinContext(ctx).Warn("响应返回", zap.Any("response_body", blw.body.String()), zap.Any("time", fmt.Sprintf("%sms", strconv.Itoa(duration))))
	statusCode := ctx.Writer.Status()
	fmt.Println(statusCode)
	//if statusCode >= 400 {
	//ok this is an request with error, let's make a record for it
	// now print body (or log in your preferred way)
	//fmt.Println("Response body: " + blw.body.String())
	//}
}
