# 开发环境搭建

## 安装 swagger api 文档生成工具 [swag](https://github.com/swaggo/swag)

```
go get -u github.com/swaggo/swag/cmd/swag
```

在 `src` 目录中执行以下命令会在 `apis` 目录中生成 api 文档

```
swag init --dir ./ --generalInfo apis/apis.go --propertyStrategy snakecase --output ./apis/docs
```

api 文档地址： <http://localhost:4869/x/apidocs/index.html>

服务启动时如果环境变量设置了 `DISABLE_GIN_SWAGGER` 会关闭 api 文档。
首次访问需经过 Basic 认证登录，登录账号密码可通过配置修改，默认为 `admin` `admin`

swag 中文文档: <https://github.com/swaggo/swag/blob/master/README_zh-CN.md>

## 使用 [air](https://github.com/cosmtrek/air) 可以根据文件变化自动编译运行服务

安装：

```
go get -u github.com/cosmtrek/air
```

在 `src` 目录中执行 `air -c .air.toml` 即可运行服务，代码修改后会自动更新 api 文档并重新编译运行

## 根据 mysql 表自动生成结构体：[table2struct](https://github.com/axiaoxin-com/table2struct)

安装：

```
go get -u github.com/axiaoxin-com/table2struct
```

## 依赖的外部 HTTP 服务的 Mock 工具： [httplive](https://github.com/gencebay/httplive)

安装：

```
go get github.com/gencebay/httplive
```

启动：

```
httplive -d `pwd`/httplive.db -p 5003
```

打开浏览器访问： `http://localhost:5003` 页面上编辑 url 和对应的返回结果保存，请求对应地址就会返回你设置的返回结果

# 配置文件

服务通过 [viper](https://github.com/spf13/viper) 加载配置文件， viper 支持的配置文件格式都可以使用。

服务启动时默认加载当前目录的 [config.default.toml](./config.default.toml) 作为配置。其中包含了服务支持的全部配置项。

服务启动时可以通过以下参数指定其他配置文件：

- `-p` 指定配置文件的所在目录
- `-c` 指定配置文件的不带格式后缀的文件名
- `-t` 指定配置文件的文件格式名

只支持从`1`个目录读取`1`个配置文件。

**建议**：在开发自己的服务时，复制当前目录的 toml 配置创建一份新的配置，再在其上进行修改或新增配置，然后通过指定参数加载自己的配置。

# 日志打印

使用 [logging](https://github.com/axiaoxin-com/logging) 的方法打印带 trace id 的日志，可通过配置文件中 `[logging]` 下的配置项进行相关设置。

配置 sentry dsn 后，`Error` 级别以上的日志会被自动采集到 Sentry 便于错误发现与定位。

# API 开发

使用 [pink-lady](http://github.com/axiaoxin-com/pink-lady) 开发 web api 服务，你只需实现 gin 的 `HandlerFunc` 并在 `apis/routes.go` 的 `Routes` 函数中注册到对应的 URL 上即可。

api 中使用 `c.Error(err)` 会将 err 保存到 context 中，打印访问日志时会以 `Error` 级别自动打印错误信息。避免同一类错误打印多次日志影响问题定位效率。

手动完整的启动服务命令：

```
go run main.go -p . -c config.default -t toml
```
