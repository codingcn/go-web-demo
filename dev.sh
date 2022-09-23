#!/usr/bin/env

go build -o ./dist/server
#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./dist/server


mkdir -p ./dist/config/
cp -r ./config/* ./dist/config/

