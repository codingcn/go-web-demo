FROM golang:1.13.4-alpine3.10

MAINTAINER Chris 5303221@gmail.com


ENV dir /data/app/go-web-demo
WORKDIR $dir
COPY . $dir
RUN  go build -mod=vendor -ldflags="-s -w" -o bin/go-web-demo

EXPOSE 9903
ENTRYPOINT [ "./bin/bin_tmp/go-web-demo", "-env", "dev" ]