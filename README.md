# pink-lady

![proj-icon](./misc/pics/logo.png)

[![go report card](https://goreportcard.com/badge/github.com/axiaoxin-com/pink-lady)](https://goreportcard.com/report/github.com/axiaoxin-com/pink-lady)
[![version-badge](https://img.shields.io/github/release/axiaoxin-com/pink-lady.svg)](https://github.com/axiaoxin-com/pink-lady/releases)
[![license](https://img.shields.io/github/license/axiaoxin-com/pink-lady.svg)](https://github.com/axiaoxin-com/pink-lady/blob/master/LICENSE)
[![issues](https://img.shields.io/github/issues/axiaoxin-com/pink-lady.svg)](https://github.com/axiaoxin-com/pink-lady/issues)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/axiaoxin-com/pink-lady/pulls)

> Pinklady is a template project of gin app, which encapsulates mysql, redis, logging, viper, swagger, middlewares and other common components.

pink-lady 是基于 Golang web 开发框架 [gin](https://github.com/gin-gonic/gin)
来进行 **API服务/WEB网站** 开发的示例项目，新建项目时可以使用它作为项目模板。

之所以叫 pink-lady 首先字面意思就是红粉佳人或粉红女郎，有这个性感的名字相信你更会好好对你的代码负责。
其次，因为 gin 就是国外六大类烈酒之一的金酒，是近百年来调制鸡尾酒时最常使用的基酒，其配方多达千种以上，
而 pink lady 是以 gin 作 base 的国标鸡尾酒之一，在这里 pink-lady 则是以 gin 作 base 的代码骨架模板之一

## 使用 pink-lady 模板创建项目

点击 <https://github.com/axiaoxin-com/pink-lady/generate> 创建你的 github 项目（使用该方式创建项目时，如需修改项目名称需手动修改）

或者手动本地创建（如想自定义项目名，推荐使用该方式）：

```
bash <(curl -s https://raw.githubusercontent.com/axiaoxin-com/pink-lady/master/misc/scripts/new_project.sh)
```

## 特性

- 使用 viper 加载配置，支持配置热更新，服务所有特性都通过配置文件控制
- 支持生成 swagger api 文档
- 封装数据库连接实例池，通过读取配置文件可以直接在代码中使用 gorm 和 sqlx 快速连接 mysql、sqlite3、postgresql、sqlserver
- 封装 redis， redis sentinel， redis cluster 连接实例池
- 封装统一的 JSON 返回结构
- 集成 sentry 搜集错误
- 内置 GinLogger 中间件打印详细的访问日志，支持不同的 http 状态码使用不同的日志级别，通过配置开关打印请求头，请求餐宿，响应体等调试信息
- 内置 GinRecovery 中间件，异常服务默认按状态码返回 JSON 错误信息，panic 错误统一交由 GinLogger 打印，支持自定义输出格式
- 使用 logging 打印日志，支持 trace id，error 以上级别自动上报到 sentry
- 支持 prometheus metrics exporter
- 支持 ratelimiter 请求限频
- 通过配置集成 go html template，可自由注册 template funcs map
- embed 静态资源编译进二进制文件中
- i18n国际化支持
- 支持类似Django的[flatpages](https://github.com/axiaoxin-com/pink-lady/blob/master/statics/flatpages/docs/README.md)
- SEO良好支持

## 使用 `pink-lady/webserver` 3 步组装一个 WEB 应用

1. 确认配置文件正确。
   配置文件必须满足能解析出指定的内容，复制或修改 [config.default.toml](https://github.com/axiaoxin-com/pink-lady/blob/master/config.default.toml) 中的配置项
2. 创建自定义中间件的 gin app `NewGinEngine` （可选）
3. 运行 web 应用服务器 `Run`。
   需传入 gin app 和在该 app 上注册 URL 路由注册函数

实现代码在`src`路径下，在 pink-lady 模板项目下，你只需关注如何实现你的业务逻辑，不用考虑如何组织项目结构和集成一些通用功能，比如数据库的连接封装，配置文件的读取，swagger 文档生成，统一的 JSON 返回结果，错误码定义，集成 Sentry 等等。

你可以在`routes`路径下实现你的 api，并在 `routes/routes.go` 的 `Routes` 函数中注册 URL 即可。外部第三方服务放在 `services` 包中进行加载或初始化。数据库模型相关定义放到 `models` 包中便于复用。

## 关于 gin

### gin 框架源码图解

![gin arch](https://github.com/axiaoxin-com/pink-lady/blob/master/misc/pics/gin_arch.svg)

### gin 中间件原理解析

<https://github.com/axiaoxin/axiaoxin/issues/17>

## 开发环境搭建

### 安装 swagger api 文档生成工具 [swag](https://github.com/swaggo/swag)

```
go get -u github.com/swaggo/swag/cmd/swag
```

在项目根目录中执行以下命令会在 `routes` 目录中生成 api 文档

```
swag init --dir ./ --generalInfo routes/routes.go --propertyStrategy snakecase --output ./routes/docs
```

api 文档地址： <http://localhost:4869/x/apidocs/index.html>

服务启动时如果环境变量设置了 `DISABLE_GIN_SWAGGER` 会关闭 api 文档。
首次访问需经过 Basic 认证登录，登录账号密码可通过配置修改，默认为 `admin` `admin`

swag 中文文档: <https://github.com/swaggo/swag/blob/master/README_zh-CN.md>

## 配置文件

服务通过 [viper](https://github.com/spf13/viper) 加载配置文件， viper 支持的配置文件格式都可以使用。

服务启动时默认加载当前目录的 `config.default.toml`

服务启动时可以通过以下参数指定其他配置文件： `-c` 指定配置文件的路径

**建议**：在开发自己的服务时，复制当前目录的 toml 配置创建一份新的配置，再在其上进行修改或新增配置，然后通过指定参数加载自己的配置。

## 日志打印

使用 [logging](https://github.com/axiaoxin-com/logging) 的方法打印带 trace id 的日志，可通过配置文件中 `[logging]` 下的配置项进行相关设置。

配置 sentry dsn 后，`Error` 级别以上的日志会被自动采集到 Sentry 便于错误发现与定位。

## API 开发

使用 [pink-lady](http://github.com/axiaoxin-com/pink-lady) 开发 web api 服务，你只需实现 gin 的 `HandlerFunc` 并在 `routes/routes.go` 的 `Routes` 函数中注册到对应的 URL 上即可。

api 中使用 `c.Error(err)` 会将 err 保存到 context 中，打印访问日志时会以 `Error` 级别自动打印错误信息。避免同一类错误打印多次日志影响问题定位效率。

手动完整的启动服务命令：

```
go run main.go -c config.default.toml
```

编译：

```
go generate
CGO_ENABLED=0 GOOS=linux go build -ldflags "-X github.com/axiaoxin-com/pink-lady/routes.BuildID=${buildid}" -o pink-lady
```

## i18n国际化支持集成方法

i18n使用gettext方式读取po文件，通过自定义脚本 `i18n.sh` 一键提取翻译文本生成pot模板，自动合并更新po文件，并使用谷歌翻译进行自动翻译。

谷歌翻译需开代理访问。

### 标记i18n字符串

1. 对golang代码中需要进行翻译的文字使用`webserver.CtxI18n(c, "文字")`、`webserver.LangI18n("en", "文字")`或`I18nString("文字")`包裹
2. 对网页模板中的翻译文字使用 `{{ _i18n .meta.Lang "文字" }}`包裹。

具体的使用示例可以参考[demo主页代码](https://github.com/axiaoxin-com/pink-lady/blob/master/routes/page_home.go)

示例代码 `./routes/page_home.go`：

```go
func PageHome(c *gin.Context) {
	meta := NewMetaData(c, webserver.CtxI18n(c, "首页"))

	data := gin.H{
		"meta":  meta,
		"alert": Alert(c, "", ""),
	}

	c.HTML(http.StatusOK, "home.html", data)
	return
}
```

其中文字"首页"将被自动翻译为多语言。

### 自动提取i18n翻译文本并自动翻译

需先安装python脚本依赖：`pip install -r ./misc/i18n/requirements.txt`

自动提取需要翻译文字生成翻译模板并完成自动谷歌翻译：`./i18n.sh`

打开对应路径（默认为`statics/i18n`）下的po文件进行校验修改，msgid为原始字符串，对应的msgstr是翻译后的语言。

### 自定义数据源提取翻译文本

如果需要从数据库等自定义数据源提取需要翻译的文本，可以通过插件机制实现：

1. 在 `misc/i18n/plugins/` 目录下创建插件文件，例如 `db_strings.py`：
```python
#!/usr/bin/env python3

from typing import Set
import os
import sys

# 添加项目根目录到 Python 路径
project_root = os.path.abspath(os.path.join(os.path.dirname(__file__), '../../..'))
sys.path.append(project_root)

# 导入你的数据库模型和连接
# from your_app.models import YourModel

def extract_strings() -> Set[str]:
    """
    从数据库提取需要翻译的字符串
    返回类型必须是 Set[str]
    """
    strings = set()
    try:
        # 实现你的数据库查询逻辑
        # 例如：
        # db = get_db_connection()
        # results = db.query(YourModel).all()
        # for item in results:
        #     strings.add(item.title)
        #     strings.add(item.description)
        pass
    except Exception as e:
        print(f"Error extracting strings from database: {e}")
    return strings
```

2. 插件文件必须实现 `extract_strings()` 函数，该函数返回一个字符串集合
3. 运行 `./i18n.sh` 时会自动加载并执行所有插件
4. 插件执行结果会被合并到最终的翻译模板中

插件系统特点：
- 模块化：每个数据源可以有自己的插件
- 可扩展：轻松添加新的数据源
- 解耦：数据源逻辑与主程序分离
- 容错：单个插件失败不影响整体运行
- 灵活：插件可以访问项目的任何模块和依赖