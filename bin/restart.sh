#!/usr/bin/env bash

# 开启错误自动退出
set -e

# 变更目录
cd /data/app/go-web-demo
# 删除正在运行的二进制文件夹
rm -f go-web-demo
# 复制发布的最新程序
cp bin_tmp/go-web-demo ./
# 修改为可执行文件
chmod a+x ./go-web-demo

PIDFILE="$PWD/go-web-demo.pid"
if [ -f $PIDFILE  ];then
    # 重启服务
    PID=$(cat $PIDFILE)
    kill -HUP $PID  || { echo "restart process failed"; exit 1; }
else
    echo "start process failed, pid file not exists"
    exit 1
fi

echo "restart success"
