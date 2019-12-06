package config

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"strings"
	"sync"
)

/**
要做的事情非常简单:
1.编写配置文件信息;
2.根据配置文件编写相应的结构体;
3.根据toml的方式读取配置文件,然后赋予这个结构体;
4.把获取的信息设置给指定的服务即可.
*/
type configStruct struct {
	Env        string `toml:"env"`
	ProjectDir string
	PIDFile    string `toml:"pid_file"`
	JWTSign    string `toml:"jwt_sign"`

	Http struct {
		Enabled bool     `toml:"enabled"`
		Listen  string   `toml:"listen"`
		URLs    []string `toml:"urls"`
	}

	Weibo struct {
		AccessToken string `toml:"access_token"`
	} `toml:"weibo"`

	Mysql struct {
		Config struct {
			Datasource string `toml:"datasource"`
			Timeout    int    `toml:"timeout"`
		} `toml:"config"`
		App struct {
			Datasource string `toml:"datasource"`
			Timeout    int    `toml:"timeout"`
		} `toml:"app"`
		Slave struct {
			App struct {
				Datasource string `toml:"datasource"`
				Timeout    int    `toml:"timeout"`
			} `toml:"app"`
			Game struct {
				Datasource string `toml:"datasource"`
				Timeout    int    `toml:"timeout"`
			} `toml:"game"`
		} `toml:"slave"`
	} `toml:"mysql"`

	Redis struct {
		Default struct {
			Addr     string `toml:"addr"`
			Password string `toml:"password"`
			DB       int    `toml:"db"`
			Timeout  int    `toml:"timeout"`
		} `toml:"default"`
	} `toml:"redis"`

	Log struct {
		LogLevel          string `toml:"log_level"`
		LogFilename       string `toml:"log_file_name"`
		MaxAge            int    `toml:"max_age"`
		MaxSize           int    `toml:"max_size"`
		MaxBackups        int    `toml:"max_backups"`
		Compress          bool   `toml:"compress"`
		LogFileNameStdout string `toml:"log_file_name_stdout"`
		LogFileNameStderr string `toml:"log_file_name_stderr"`
	} `toml:"log"`
}

var C *configStruct

// go run /data/go/src/gameserver/main.go  -project_dir=/data/go/src/gameserver -env=dev
// go run /data/go/src/gameserver/main.go  -config=/data/go/src/gameserver/config/dev.toml

var (
	once       sync.Once
	projectDir = flag.String("project_dir", "./", "Location of the project_dir.")
	configFile = flag.String("config", "", "specify configuration file path")
	env        = flag.String("env", "local", "Location of the env.")
)

func init() {
	flag.Parse()
	once.Do(initConfig)

}
func initConfig() {

	projectDirFormat := strings.TrimRight(*projectDir, "/")
	configFilename := projectDirFormat + "/config/local.toml"
	if *configFile != "" {
		configFilename = *configFile
	} else if *env == "local" {
		configFilename = projectDirFormat + "/config/local.toml"
	} else if *env == "dev" {
		configFilename = projectDirFormat + "/config/dev.toml"
	} else if *env == "prod" {
		configFilename = projectDirFormat + "/config/prod.toml"
	} else {
		log.Fatalln("启动错误：env参数值必须为dev/prod/local")
	}
	fmt.Println("load config file:" + configFilename)
	if _, err := toml.DecodeFile(configFilename, &C); err != nil {
		log.Fatalln(err)
	}
	C.ProjectDir = projectDirFormat
}
