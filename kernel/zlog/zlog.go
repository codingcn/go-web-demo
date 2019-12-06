package zlog

import (
	"github.com/gin-gonic/gin"
	"go-web-demo/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strconv"
	"sync"
	"time"
)

const loggerKey = iota

type logger struct {
	*zap.Logger
}

var (
	once   sync.Once
	Logger *logger // 继承*zap.Logger并添加一些自定义方法
)

func init() {
	once.Do(initLogConfig)
}

func initLogConfig() {
	// 日志地址 "out.log" 自定义
	lp := config.C.Log.LogFilename
	// 日志级别 DEBUG,ERROR, INFO
	lv := config.C.Log.LogLevel
	var level zapcore.Level
	//debug<info<warn<error<fatal<panic
	switch lv {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	hook := lumberjack.Logger{
		Filename:   lp,                      // 日志文件路径
		MaxSize:    config.C.Log.MaxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: config.C.Log.MaxBackups, // 日志文件最多保存多少个备份
		MaxAge:     config.C.Log.MaxAge,     // 文件最多保存多少天
		Compress:   config.C.Log.Compress,   // 是否压缩
	}
	// 是否 DEBUG
	if config.C.Env != "prod" {
		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
				TimeKey:        "ts",
				LevelKey:       "level",
				NameKey:        "Logger",
				CallerKey:      "caller",
				MessageKey:     "msg",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
				EncodeTime:     timeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.FullCallerEncoder,
			}), // 编码器配置
			//zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),      // 打印到控制台
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
			level,                                                                           // 日志级别
		)

		// 开启开发模式，堆栈跟踪
		caller := zap.AddCaller()
		// 开启文件及行号
		development := zap.Development()
		// 设置初始化字段
		//filed := zap.Fields(zap.String("serviceName", "serviceName"))
		// 构造日志
		Logger = &logger{zap.New(core, caller, development, zap.AddStacktrace(zap.ErrorLevel))}

	} else {
		encoderConfig := zap.NewProductionEncoderConfig()
		Logger = &logger{zap.New(zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置(生产环境使用json)
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
			level,                                                                           // 日志级别
		), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))}
	}
}

// 自定义时间编码器
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
	//enc.AppendString(t.Format("2006-01-02 15:04:05.000000000"))
}

// 给指定的context添加字段
func (l *logger)NewContext(ctx *gin.Context, fields ...zapcore.Field) {
	ctx.Set(strconv.Itoa(loggerKey), l.WithContext(ctx).With(fields...))
}

// 从指定的context返回一个zap实例
func (l *logger)WithContext(ctx *gin.Context) *logger {
	if ctx == nil {
		return Logger
	}
	zl, _ := ctx.Get(strconv.Itoa(loggerKey))
	ctxLogger, ok := zl.(*zap.Logger)
	if ok {
		return &logger{ctxLogger}
	}
	return Logger
}
