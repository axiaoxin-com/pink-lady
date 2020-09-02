package webserver

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

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

// GinRecovery gin recovery 中间件
// save err in context and abort with 500
func GinRecovery(statusHandler ...func(c *gin.Context, status int, data interface{}, err error, extraMsgs ...interface{})) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			status := c.Writer.Status()

			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				if brokenPipe {
					// save err in context
					c.Error(errors.New(fmt.Sprint("Broken pipe:", err, "\n", string(debug.Stack()))))
					if len(statusHandler) > 0 {
						status = http.StatusInternalServerError
						statusHandler[0](c, status, nil, errors.New(http.StatusText(status)))
					} else {
						c.AbortWithStatus(http.StatusInternalServerError)
					}
					return
				}

				// save err in context
				c.Error(errors.New(fmt.Sprint("Recovery from panic:", err, "\n", string(debug.Stack()))))
				if len(statusHandler) > 0 {
					status = http.StatusInternalServerError
					statusHandler[0](c, status, nil, errors.New(http.StatusText(status)))
				} else {
					c.AbortWithStatus(http.StatusInternalServerError)
				}
				return
			}

			if len(statusHandler) > 0 && status >= 400 {
				statusHandler[0](c, status, nil, errors.New(http.StatusText(status)))
			}
		}()

		c.Next()
	}
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
		GinRecovery(response.Respond),
	}
	return m
}
