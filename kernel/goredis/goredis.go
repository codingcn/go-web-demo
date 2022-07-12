package goredis

import (
	"github.com/go-redis/redis"
	"go-web-demo/kernel/tconfig"
	"go-web-demo/kernel/zlog"
	"time"
)

type connect struct {
	Default *redis.Client
}

var Connect connect

type Tx = redis.Tx

var Nil = redis.Nil

func Init() {
	Connect.Default = redis.NewClient(&redis.Options{
		Addr:         tconfig.C.GetString("redis.default.addr"),
		Password:     tconfig.C.GetString("redis.default.password"),
		DB:           tconfig.C.GetInt("redis.default.db"),
		PoolSize:     20,
		PoolTimeout:  2 * time.Minute,
		IdleTimeout:  10 * time.Minute,
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	})
	p, err := Connect.Default.Ping().Result()
	if err != nil {
		panic(err)
	}
	zlog.Logger.Sugar().Info("goredis ", tconfig.C.GetString("redis.default.addr"), " is ", p)
}
