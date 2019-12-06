package gorm

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"go-web-demo/config"
	"time"
)

type connect struct {
	Default    *gorm.DB
}

var Connect connect

func InitDbConnect() {
	var (
		err error
	)
	// 公共数据库连接列表

	Connect.Default, err = gorm.Open("mysql", config.C.Mysql.App.Datasource)
	Connect.Default.LogMode(true)
	Connect.Default.DB().SetMaxOpenConns(10)
	Connect.Default.DB().SetMaxIdleConns(10)
	Connect.Default.DB().SetConnMaxLifetime(1 * time.Hour)
	if err != nil {
		panic(err)
	}
	if err := Connect.Default.DB().Ping(); err != nil {
		panic(err)
	}

}
