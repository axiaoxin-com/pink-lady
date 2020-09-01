package webserver

import (
	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/pink-lady/response"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	// GinPprofURLPath 设置 gin 中的 pprof url 注册路径，可以通过外部修改
	GinPprofURLPath = "/x/pprof"
)

// NewGinEngine 根据参数创建 gin 的 router engine
// middlewares 需要使用到的中间件列表，默认不为 engine 添加任何中间件
func NewGinEngine(middlewares ...gin.HandlerFunc) *gin.Engine {
	// set gin mode
	gin.SetMode(viper.GetString("server.mode"))

	engine := gin.New()

	// use middlewares
	for _, middleware := range middlewares {
		engine.Use(middleware)
	}

	if viper.GetBool("server.pprof") {
		pprof.Register(engine, GinPprofURLPath)
	}
	return engine
}

// GinBasicAuth gin 的基础认证中间件
// 加到 gin app 的路由中可以对该路由添加 basic auth 登录验证
// 传入 username 和 password 对可以替换默认的 username 和 password
func GinBasicAuth(args ...string) gin.HandlerFunc {
	username := viper.GetString("basic_auth.username")
	password := viper.GetString("basic_auth.password")
	if len(args) == 2 {
		username = args[0]
		password = args[1]
	}
	logging.Debug(nil, "Basic auth username:"+username+" password:"+password)
	return gin.BasicAuth(gin.Accounts{
		username: password,
	})
}

// DefaultGinMiddlewares 默认的 gin server 使用的中间件列表
func DefaultGinMiddlewares() []gin.HandlerFunc {
	m := []gin.HandlerFunc{
		// 记录请求处理日志，最顶层执行
		logging.GinLoggerWithConfig(logging.GinLoggerConfig{
			SkipPaths:           viper.GetStringSlice("logging.access_logger.skip_paths"),
			EnableDetails:       viper.GetBool("logging.access_logger.enable_details"),
			EnableContextKeys:   viper.GetBool("logging.access_logger.enable_context_keys"),
			EnableRequestHeader: viper.GetBool("logging.access_logger.enable_request_header"),
			EnableRequestForm:   viper.GetBool("logging.access_logger.enable_request_form"),
			EnableRequestBody:   viper.GetBool("logging.access_logger.enable_request_body"),
			EnableResponseBody:  viper.GetBool("logging.access_logger.enable_response_body"),
		}),
		// 捕获 panic 保存到 context 中由 GinLogger 统一打印， panic 时返回 500 JSON
		logging.GinRecovery(response.Respond),
	}
	return m
}
