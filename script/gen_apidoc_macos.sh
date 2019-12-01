#! /usr/bin/env bash
# 生成swag api文档
PROJECT_PATH=$(dirname $(dirname $(greadlink -f $0)))
MAIN_PATH=${PROJECT_PATH}/app/
# swag init必须在main.go所在的目录下执行，否则必须用--dir参数指定main.go的路径
swag init --dir ${MAIN_PATH} --generalInfo apis/apis.go --propertyStrategy camelcase --output ${MAIN_PATH}/apis/docs
