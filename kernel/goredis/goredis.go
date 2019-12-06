package goredis

import (
	"github.com/go-redis/redis"
	"go-web-demo/config"
	"go-web-demo/kernel/zlog"
	"time"
)

type connect struct {
	Default *redis.Client
}

var Connect connect

type Tx = redis.Tx

var Nil = redis.Nil

func InitConnect() {

	Connect.Default = redis.NewClient(&redis.Options{
		Addr:         config.C.Redis.Default.Addr,
		Password:     config.C.Redis.Default.Password,
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
	zlog.Logger.Sugar().Info("goredis ", config.C.Redis.Default.Addr, " is ", p)
}
