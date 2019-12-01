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

pink-lady 是基于 Golang web 开发框架 [gin](https://github.com/gin-gonic/gin)
来进行 HTTP API 开发的示例项目，新建项目时可以使用它作为项目模板。

之所以叫 pink-lady 首先字面意思就是红粉佳人或则粉红女郎，有这个性感的名字相信你更会好好对你的代码负责。
其次，因为 gin 就是国外六大类烈酒之一的金酒，是近百年来调制鸡尾酒时最常使用的基酒，其配方多达千种以上，
而 pink lady 是以 gin 作 base 的国标鸡尾酒之一，在这里 pink-lady 则是以 gin 作 base 的代码骨架模板之一

## 依赖管理

go mod

## 配置

使用 [viper](https://github.com/spf13/viper) 读取配置，
配置文件参考`app/config.toml.example`生成你自己的`app/config.toml`文件，
可以使用任意viper支持的文件格式，推荐使用 [toml](https://github.com/toml-lang/toml)

代码中读取配置在配置文件中新增你的配置项，然后直接使用`viper.GetXXX("a.b.c")`即可读取

## 接口实现

在`app/apis/routes.go`中添加URL并指定你的handleFunc，handleFunc推荐在`app/apis`下按业务模块新建文件的形式组织
可复用的代码可以在`app/services`下以目录或者文件的形式按需组织

API版本号定义在`app/api/apis.go`中，可以手动修改值，但不要修改代码格式，自动生成API文档依赖这个格式。

## Demo

提供了一个demo接口实现用于参考，涉及

- `app/apis/routes.go`
- `app/apis/demo*`
- `app/models/demomod/`
- `app/services/demosvc/`

## 访问DB

使用 [gorm](https://github.com/jinzhu/gorm) 访问db，在`app/models`中定义数据库模型，使用`app/db`包获取db实例，数据库实例按配置文件中的配置全部生成。
使用配置中的`instance`的值可以获取对应数据库实例，例如获取MySQL配置中的`instance = "default"`的数据库实例使用`db.MySQL("default")`即可，其他实例类似。

建议针对你的DB变更做迁移备份，这里推荐使用 [goose](https://github.com/pressly/goose)

## 日志

使用 [zap](https://github.com/uber-go/zap) 打印日志，普通日志直接使用全局的`logging.Logger`的方式打印，打印带有context中requestid的日志使用`logging.CtxLogger(c)`

## 中间件

中间件存放在`app/middleware`中，在`app/router/router.go`进行注册。

## 接口返回

接口返回统一的 JSON 结构

## Swagger API文档

使用 [swag](https://github.com/swaggo/swag) 生成api文档，
运行`misc/gen_apidoc.sh`可以根据swag支持的注释格式生成swagger api文档

需要先安装swag：

```
go get -u -v github.com/swaggo/swag/cmd/swag
```

生成的文档在`app/apis/docs`中，你自己的文档不要放在这里。

生成文档后启动服务，访问 <http://127.0.0.1:4869/x/apidocs/index.html> 即可看到效果。登录账号：`admin`  密码：`!admin`
浏览器请求地址host必须和`app/apis/apis.go`中的`// @host`指定的注释一致，如果更改请求地址需要同步修改这里。

## 快速创建以 pink-lady 为模板的项目

```
bash <(curl -s https://raw.githubusercontent.com/axiaoxin/pink-lady/master/misc/new-project.sh)
```
