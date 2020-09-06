# pink-lady

![proj-icon](./misc/pics/logo.png)

[![Build Status](https://travis-ci.org/axiaoxin-com/pink-lady.svg?branch=master)](https://travis-ci.org/axiaoxin-com/pink-lady)
[![go report card](https://goreportcard.com/badge/github.com/axiaoxin-com/pink-lady)](https://goreportcard.com/report/github.com/axiaoxin-com/pink-lady)
[![version-badge](https://img.shields.io/github/release/axiaoxin-com/pink-lady.svg)](https://github.com/axiaoxin-com/pink-lady/releases)
[![license](https://img.shields.io/github/license/axiaoxin-com/pink-lady.svg)](https://github.com/axiaoxin-com/pink-lady/blob/master/LICENSE)
[![issues](https://img.shields.io/github/issues/axiaoxin-com/pink-lady.svg)](https://github.com/axiaoxin-com/pink-lady/issues)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/axiaoxin-com/pink-lady/pulls)

> Pinklady is a template project of gin app, which encapsulates mysql, redis, logging, viper, swagger, middlewares and other common components.

pink-lady 是基于 Golang web 开发框架 [gin](https://github.com/gin-gonic/gin)
来进行 HTTP API 开发的示例项目，新建项目时可以使用它作为项目模板。

之所以叫 pink-lady 首先字面意思就是红粉佳人或则粉红女郎，有这个性感的名字相信你更会好好对你的代码负责。
其次，因为 gin 就是国外六大类烈酒之一的金酒，是近百年来调制鸡尾酒时最常使用的基酒，其配方多达千种以上，
而 pink lady 是以 gin 作 base 的国标鸡尾酒之一，在这里 pink-lady 则是以 gin 作 base 的代码骨架模板之一

## 特性

- 使用 viper 加载配置，支持配置热更新，服务所有特性都通过配置文件控制
- 支持生成 swagger api 文档
- 封装数据库连接实例池，通过读取配置文件可以直接在代码中使用 gorm 和 sqlx 快速连接 mysql、sqlite3、postgresql、sqlserver
- 封装 redis， redis sentinel， redis cluster 连接实例池
- 封装统一的 JSON 返回结构
- 集成 sentry 搜集错误
- 内置 GinLogger 中间件打印详细的访问日志，支持不同的 http 状态码使用不同的日志级别，通过配置开关打印请求头，请求餐宿，响应体等调试信息
- 内置 GinRecovery 中间件，异常服务默认按状态码返回 JSON 错误信息，panic 错误统一交由 GinLogger 打印，支持自定义输出格式
- 内置 GinTimeout 中间件，可以为请求处理设置超时时间，超时时间到达后返回 503 JSON 错误信息
- 使用 logging 打印日志，支持 trace id，error 以上级别自动上报到 sentry
- 支持 prometheus metrics exporter

## 使用 `pink-lady/webserver` 3 步组装一个 WEB 应用

1. 确认配置文件正确。
   配置文件必须满足能解析出指定的内容，复制或修改 [config.default.toml](./src/config.default.toml) 中的配置项
2. 创建自定义中间件的 gin app `NewGinEngine` （可选）
3. 运行 web 应用服务器 `Run`。
   需传入 gin app 和在该 app 上注册 URL 路由注册函数

实现代码在`src`路径下，在 pink-lady 模板项目下，你只需关注如何实现你的业务逻辑，不用考虑如何组织项目结构和集成一些通用功能，比如数据库的连接封装，配置文件的读取，swagger 文档生成，统一的 JSON 返回结果，错误码定义，集成 Sentry 等等。

你可以在`apis`路径下实现你的 api，并在 `apis/routes.go` 的 `Routes` 函数中注册 URL 即可。可复用的业务代码可以放到 `handlers` 包中，便于比如定时任务复用业务逻辑代码。数据库模型相关定义放到 `models` 包中便于复用。

## 关于 gin

### gin 框架源码图解

![gin arch](./misc/pics/gin_arch.svg)

### gin 中间件原理解析

<https://github.com/axiaoxin/axiaoxin/issues/17>
