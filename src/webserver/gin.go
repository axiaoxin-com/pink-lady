package webserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

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
// save err in context and abort with recoveryHandler
func GinRecovery(recoveryHandler ...func(c *gin.Context, status int, data interface{}, err error, extraMsgs ...interface{})) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			status := c.Writer.Status()
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				status = http.StatusInternalServerError
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
					c.Error(fmt.Errorf("Broken pipe: %v\n%s", err, string(debug.Stack())))
					if len(recoveryHandler) > 0 {
						c.Abort()
						recoveryHandler[0](c, status, nil, errors.New(http.StatusText(status)))
						return
					}
					c.AbortWithStatus(status)
					return
				}

				// save err in context when panic
				c.Error(fmt.Errorf("Recovery from panic: %v\n%s", err, string(debug.Stack())))
			}

			// 状态码大于 400 的以 status handler 进行响应
			if status >= 400 {
				// 有 handler 时使用 handler 返回
				if len(recoveryHandler) > 0 {
					c.Abort()
					recoveryHandler[0](c, status, nil, errors.New(http.StatusText(status)))
					return
				}
				// 否则只返回状态码
				c.AbortWithStatus(status)
				return
			}
		}()
		c.Next()
	}
}

// GinTimeout 设置 gin handler 请求处理超时时间，超时时间到达直接返回 503
func GinTimeout(timeoutSec ...int) gin.HandlerFunc {
	// 使用 goroutine 运行 Next() 处理请求，同时 select 监听 timeout context 的超时
	// 由于在 goroutine 中发生 panic 会导致进程退出，所以需要在 goroutine 中进行 recover 避免进程退出
	// 使用 channel 的方式将 goroutine 中的 panic 抛到外层 panic ，统一由 GinRecovery 中间件处理

	timeout := viper.GetDuration("server.handler_timeout") * time.Second
	if len(timeoutSec) > 0 {
		timeout = time.Duration(timeoutSec[0]) * time.Second
	}
	return func(c *gin.Context) {
		// 设置超时 context
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		c.Request.WithContext(ctx)

		// 处理请求
		// 超时发生时，由于先返回了超时的信息，而 Next 仍然在继续执行
		// 当 goroutine 的 Next() 执行完成后返回结果时会 panic -> http: wrote more than the declared Content-Length
		// 但因为在发生超时逻辑已进入 select 之后的逻辑，此时发生 panic 只会保存到 panicChan 中而不会被抛到上层
		done := make(chan struct{})
		panicChan := make(chan interface{}, 1)
		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			c.Next()
			close(done) // 关闭 done channel 通知 select 任务已完成
		}()

		// 阻塞监听 channel
		select {
		case p := <-panicChan:
			// 上抛异常
			panic(p)
		case <-ctx.Done():
			// 超时返回 503
			c.AbortWithStatus(http.StatusServiceUnavailable)
		case <-done:
			// 正常完成处理无需其他操作
		}
	}
}

// DefaultGinMiddlewares 默认的 gin server 使用的中间件列表
func DefaultGinMiddlewares() []gin.HandlerFunc {
	m := []gin.HandlerFunc{
		// 记录请求处理日志，最顶层执行
		logging.GinLoggerWithConfig(logging.GinLoggerConfig{
			SkipPaths:           viper.GetStringSlice("logging.access_logger.skip_paths"),
			SkipPathRegexps:     viper.GetStringSlice("logging.access_logger.skip_path_regexps"),
			EnableDetails:       viper.GetBool("logging.access_logger.enable_details"),
			EnableContextKeys:   viper.GetBool("logging.access_logger.enable_context_keys"),
			EnableRequestHeader: viper.GetBool("logging.access_logger.enable_request_header"),
			EnableRequestForm:   viper.GetBool("logging.access_logger.enable_request_form"),
			EnableRequestBody:   viper.GetBool("logging.access_logger.enable_request_body"),
			EnableResponseBody:  viper.GetBool("logging.access_logger.enable_response_body"),
		}),
		// 捕获 panic 保存到 context 中由 GinLogger 统一打印， panic 时返回 500 JSON
		GinRecovery(response.Respond),
		GinTimeout(),
	}
	return m
}
