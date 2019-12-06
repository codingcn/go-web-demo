#!/usr/bin/env bash

# 需要配置crontab脚本
# * * * * * flock -xno /tmp/go-web-demo-wathdog.lock -c 'sh /data/app/go-web-demo/watchdog.sh >> /data/app/go-web-demo/logs/watchdog.zlog 2>&1'
# 必须添加环境变量APP_ENV，可选值dev|prod


# 开启错误自动退出
set -o errexit
# 执行命令的目录
baseDir='/data/app/go-web-demo'
# hummer_ipay 启动命令
# startCMD='./go-web-demo --env=prod 1> ./logs/stdout.zlog 2> ./logs/stderr.zlog'
startCMD="./go-web-demo --env=${APP_ENV}"
pid='go-web-demo.pid'

# 脚本必须以root身份运行
user=`whoami`
if [ "$user" != "root" ]; then
    echo "this tool must run as *root*"
    exit 1
fi

# 切换到执行目录
cd $baseDir

# while true;do
for((i=1;i<20;i++));do
    # 检查服务，如果没有启动，则执行启动命令
    now=`date '+%Y-%m-%d %H:%M:%S'`
    ret=`ps aux | grep "$startCMD" | grep -v grep | wc -l`
    if [ $ret -eq 0 ]; then
        if [ -e $pid ];then
            rm -f $pid || {echo "remove pid file failed"}
        fi
        echo "$now process not exists ,start process now... "
        if eval $startCMD;then
            echo "$now start process done!"
        else
            echo "$now start process failed"
        fi
    else
        echo "$now process is running..."
    fi
    sleep 3
done
