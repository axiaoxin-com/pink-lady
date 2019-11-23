# pink-lady

![proj-icon](https://raw.githubusercontent.com/axiaoxin/pink-lady/master/misc/pinklady.png)

[![Build Status](https://travis-ci.org/axiaoxin/pink-lady.svg?branch=master)](https://travis-ci.org/axiaoxin/pink-lady)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/b906dd1655074f60bf93a7c592d29204)](https://www.codacy.com/app/axiaoxin/pink-lady?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=axiaoxin/pink-lady&amp;utm_campaign=Badge_Grade)
[![go report card](https://goreportcard.com/badge/github.com/axiaoxin/pink-lady)](https://goreportcard.com/report/github.com/axiaoxin/pink-lady)
[![codecov](https://codecov.io/gh/axiaoxin/pink-lady/branch/master/graph/badge.svg)](https://codecov.io/gh/axiaoxin/pink-lady)
![version-badge](https://img.shields.io/github/release/axiaoxin/pink-lady.svg)
![downloads](https://img.shields.io/github/downloads/axiaoxin/pink-lady/total.svg)
![license](https://img.shields.io/github/license/axiaoxin/pink-lady.svg)
![issues](https://img.shields.io/github/issues/axiaoxin/pink-lady.svg)
[![Say Thanks!](https://img.shields.io/badge/Say%20Thanks-!-1EAEDB.svg)](https://saythanks.io/to/axiaoxin)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/axiaoxin/pink-lady/pulls)

pink-lady是使用[gin](https://github.com/gin-gonic/gin)做web开发的demo项目，你可以通过执行以下以下命令复制一份在这个基础上来开始你自己新的web项目。

    # macos
    bash <(curl -s https://raw.githubusercontent.com/axiaoxin/pink-lady/master/misc/new-project.macos.sh)
    # linux
    bash <(curl -s https://raw.githubusercontent.com/axiaoxin/pink-lady/master/misc/new-project.linux.sh)

所有代码在app目录中，misc中存放各类脚本，在go 1.13测试通过。


# 如何开始开发

## 配置
配置文件参考`app/config.toml.example`生成你自己的`app/config.toml`文件，使用[viper](https://github.com/spf13/viper)读取
配置文件可以使用任意viper支持的文件格式，推荐使用[toml](https://github.com/toml-lang/toml)

代码中读取配置在配置文件中新增你的配置项，然后直接使用`viper.GetXXX("a.b.c")`即可读取

## 接口实现
在`app/apis/routes.go`中添加URL并指定你的handleFunc，handleFunc可以在`app/apis`下以目录或者文件的形式自己按实际情况组织
可复用的代码可以在`app/services`下以目录或者文件的形式按需组织

API版本号定义在`app/api/init.go`中，可以手动修改值，但不要修改代码格式，自动生成API文档依赖这个格式。

## 访问DB
使用(gorm](https://github.com/jinzhu/gorm)访问db，在`app/models`中定义数据库模型，使用`app/db`包获取db实例，db实例按配置文件中的配置全部生成。
使用配置中的instance的值可以获取对应db实例，例如获取MySQL配置中的`instance = "default"`的数据库实例使用`db.MySQL("default")`即可，其他db实例类似。

## 日志
使用[logrus](https://github.com/sirupsen/logrus)打印日志，日志不打印到文件全部输出到标准输出。
普通日志直接使用logrus的方式答应，打印请求信息使用全局的logger`utils.Logger`进行普通日志打印，该Logger是一个logrus的Entry实例，用法直接参考logrus

如果需要自动打印RequestID，必须使用`utils.CtxLogger(ctx)`实时获取带RequestID的Logger

## 中间件
中间件存放在`app/middleware`中，在`app/router/router.go`进行注册。

## 返回值
使用`app/response`中的方法返回统一结构的json，返回码参数为`RetCode`结构体对象，在`app/retcode`中新增返回码。一个返回码结构体对象包含实际的code和code对应的message，同一个code要对应不同message必须新增RetCode结构体对象

## Swagger API文档
使用[swag](https://github.com/swaggo/swag)生成api文档，运行`misc/gen_apidoc_xxx.sh`可以根据swag支持的注释格式生成swagger api文档

需要先安装swag：

    go get -u -vgithub.com/swaggo/swag/cmd/swag

文档相关的代码在`app/docs`中，你自己的文档不要放在这里。

生成文档后启动服务，访问<http://127.0.0.1:4869/x/apidocs/index.html>即可看到效果。
浏览器请求地址host必须和`app/apis/init.go`中的`// @host`指定的注释一致，如果更改请求地址需要同步修改这里。
