package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	// DefaultAPIDocsUsername apidocs 默认用户名
	DefaultAPIDocsUsername = "admin"
	// DefaultAPIDocsPassword apidocs 默认密码
	DefaultAPIDocsPassword = "!admin"
)

// BasicAuth use username and password in config to auth
func BasicAuth() gin.HandlerFunc {
	username := viper.GetString("apidocs.basicauth.username")
	if username == "" {
		username = DefaultAPIDocsUsername
	}
	password := viper.GetString("apidocs.basicauth.password")
	if password == "" {
		password = DefaultAPIDocsPassword
	}
	return gin.BasicAuth(gin.Accounts{
		username: password,
	})
}
