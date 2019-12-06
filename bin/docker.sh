#!/usr/bin/env bash

tag="1.0.3"
# kubectl -n test-ns set image deployments/go-web-demo go-web-demo=docker.xxxxx.cn:8066/micro-services/go-web-demo
# https://www.imooc.com/article/48682
docker build  -t docker.xxxxx.cn:8066/micro-services/go-web-demo:${tag} .
docker push docker.xxxxx.cn:8066/micro-services/go-web-demo:${tag}

#docker stop docker.xxxxx.cn:8066/micro-services/go-web-demo
#docker rm docker.xxxxx.cn:8066/micro-services/go-web-demo:${tag}
#docker run --name go-web-demo -p 9903:9903 -d  docker.xxxxx.cn:8066/micro-services/go-web-demo:${tag}


# 请求自动部署脚本