package torm

import (
	"go-web-demo/kernel/tconfig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

type connect struct {
	Default *gorm.DB
	once    sync.Once
}

var Connect connect

func Init() {
	Connect.once.Do(func() {
		initConnect()
	})
}
func initConnect() {
	// 公共数据库连接列表
	var err error
	Connect.Default, err = gorm.Open(mysql.Open(tconfig.C.GetString("mysql.default.datasource")), &gorm.Config{
		ConnPool: nil,
		//Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

}
