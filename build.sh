#! /usr/bin/env bash

buildid=`TZ=Asia/Shanghai date +'%Y%m%d'`
# 二进制编译
echo ">> build go"
CGO_ENABLED=0 GOOS=linux go build -ldflags "-X github.com/axiaoxin-com/pink-lady/routes.BuildID=${buildid}" -o pink-lady

#tar czvf pink-lady.tar.gz pink-lady config.default.toml statics
