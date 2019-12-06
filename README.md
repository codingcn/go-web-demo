## Golang web框架demo

**Tips：该项目不包含具体业务逻辑，只提供一个项目空架子**

> 几个知名的Go语言Web框架(Echo、Gin和Buffalo)由于没有完备支持所有功能，并不能算是真正意义上的Web框架，但大部分go社区认为它们是的;


这类框架虽然带来了更多的灵活性，却导致了大多数`Golang`项目结构十分杂乱(不仅限于web项目...)；

团队开发中，没有制定严格的编码规范，同事的业务代码经常随意乱扔，

**该项目主要提供了一个直接可用的空壳架子，但封装了许多优秀的组件应用组件，可以参考快速构建一个完整的项目。**



## 功能特性

* 项目结构借鉴[laravel/laravel](https://github.com/laravel/laravel.git)风格，基于gin构建；
* 平滑更新支持 [cloudflare/tableflip](https://github.com/cloudflare/tableflip)；
* 请求验证器及中间件支持： 跨域； prometheus； [dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go)；
* 日志链路及日志轮转支持： [zap](https://go.uber.org/zap) +context+middleware+[lumberjack.v2](https://gopkg.in/natefinch/lumberjack.v2)；
* 数据库支持：[gorm](https://github.com/jinzhu/gorm)+[go-redis](https://github.com/go-redis/redis)；
* 静态配置支持：[toml](https://github.com/BurntSushi/toml)；
* 计划任务支持：[robfig/cron](https://github.com/robfig/cron)；


## 运行

启动 项目
```
git clone https://github.com/codingcn/go-web-demo.git
cd go-web-demo
go run main.go -env=dev
```

平滑升级（替换go可执行文件后，发送指定信号量）
```
kill -HUP  `cat go-web-demo.pid`
```


