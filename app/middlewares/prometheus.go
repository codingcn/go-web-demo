package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"strings"
	"time"
)

var (
	//HTTPReqDuration metric:http_request_duration_seconds
	HTTPReqDuration *prometheus.HistogramVec
	//HTTPReqTotal metric:http_request_total
	HTTPReqTotal *prometheus.CounterVec
)

func init() {
	// 监控接口请求耗时
	// HistogramVec 是一组Histogram
	HTTPReqDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "The HTTP request latencies in seconds.",
		Buckets: nil,
	}, []string{"method", "path"})
	// 这里的"method"、"path" 都是label
	// 监控接口请求次数
	// HistogramVec 是一组Histogram
	HTTPReqTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests made.",
	}, []string{"method", "path", "status"})
	// 这里的"method"、"path"、"status" 都是label
	prometheus.MustRegister(
		HTTPReqDuration,
		HTTPReqTotal,
	)
}

// /api/epgInfo/1371648200  -> /api/epgInfo
func parsePath(path string) string {
	itemList := strings.Split(path, "/")
	if len(path) >= 4 {
		return strings.Join(itemList[0:3], "/")
	}
	return path
}

//Metric metric middleware
func Metric() gin.HandlerFunc {
	return func(c *gin.Context) {
		tBegin := time.Now()
		duration := float64(time.Since(tBegin)) / float64(time.Second)

		//path := parsePath(c.Request.URL.Path)
		path := c.Request.URL.Path

		// 请求数加1
		HTTPReqTotal.With(prometheus.Labels{
			"method": c.Request.Method,
			"path":   path,
			"status": strconv.Itoa(c.Writer.Status()),
		}).Inc()

		//  记录本次请求处理时间
		HTTPReqDuration.With(prometheus.Labels{
			"method": c.Request.Method,
			"path":   path,
		}).Observe(duration)
		c.Next()
	}
}
