package tconfig

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"sync"
)

var (
	once    sync.Once
	AppName string
	AppEnv  string
	C       *viper.Viper
)

func Init() {
	once.Do(func() {
		initConfig()
	})
}

func initConfig() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env2 := flag.String("env", "local", "环境参数")
		flag.Parse()
		env = *env2

	}
	if env == "" {
		env = "local"
	}
	AppEnv = env
	fmt.Println("当前配置环境：", env)
	getConfig("./config/" + env + ".toml")

}
func getConfig(path string) {

	//viper.RemoteConfig = &Config{}
	//C = viper.New()
	//endpoint := "http://127.0.0.1:2379"
	//C.AddRemoteProvider("etcd", endpoint, path)
	//C.SetConfigType("toml")
	//err := C.ReadRemoteConfig()
	//if err != nil {
	//	log.Fatal("ReadRemoteConfig error", err, path)
	//}
	//C.WatchRemoteConfigOnChannel()
	//

	C = viper.New()
	C.SetConfigFile(path)
	err := C.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
