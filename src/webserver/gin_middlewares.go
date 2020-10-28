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
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

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
