package middleware

import (
	"pink-lady/app/utils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func TestBasicAuth(t *testing.T) {
	r := gin.New()
	viper.SetDefault("apidocs.basicauth.username", "admin")
	viper.SetDefault("apidocs.basicauth.password", "admin")
	r.Use(BasicAuth())
	r.GET("/", func(c *gin.Context) { return })
	respRecorder := utils.PerformRequest(r, "GET", "/", nil)
	if respRecorder.Code != 401 {
		t.Fatal("BasicAuth not work")
	}
}
