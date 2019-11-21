#! /usr/bin/env bash
PROJECT_PATH=$(dirname $(dirname $(greadlink -f $0)))
APP_PATH=${PROJECT_PATH}/app
cd ${APP_PATH}
swag init -g apis/init.go
cd -
