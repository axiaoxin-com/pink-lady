package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// BasicAuth use username and password in config to auth
func BasicAuth() gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		viper.GetString("apidocs.basicauth.username"): viper.GetString("apidocs.basicauth.password"),
	})
}
