# pink-lady

![proj-icon](./misc/logo.png)

[![Build Status](https://travis-ci.org/axiaoxin/pink-lady.svg?branch=master)](https://travis-ci.org/axiaoxin/pink-lady)
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

实现代码在`src`路径下，在 pink-lady 模板项目下，你只需关注如何实现你的业务逻辑，不用考虑如何组织项目结构和集成一些通用功能，比如数据库的连接封装，配置文件的读取，swagger 文档生成，统一的 JSON 返回结果，错误码定义，集成 Sentry 等等。

你可以在`apis`路径下实现你的 api，并在 `apis/apis.go` 的 `Routes` 函数中注册 URL 即可。可复用的业务代码可以放到 `handler` 包中，便于比如定时任务复用业务逻辑代码。数据库模型相关定义放到 `models` 包中便于复用。

## 使用 `pink-lady/webserver` 3 步组装一个 WEB 应用

1. 初始化 web 应用的配置信息到 viper `InitConfig` 。
   配置文件必须满足能解析出指定的内容，参考 [config.default.toml](./src/config.default.toml) 中的配置项
2. 创建自定义中间件的 gin app `NewGinEngine`
3. 运行 web 应用服务器 `Run`。
   需传入 gin app 和在该 app 上注册 URL 路由注册函数
