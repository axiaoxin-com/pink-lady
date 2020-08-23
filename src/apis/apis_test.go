package apis

import (
	"testing"

	"github.com/axiaoxin-com/goutils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func TestRegisterRoutes(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	// Register 中的 basic auth 依赖 viper 配置
	viper.Set("basic_auth.username", "admin")
	viper.Set("basic_auth.password", "admin")
	defer viper.Reset()

	Register(r)
	recorder, err := goutils.RequestHTTPHandler(r, "GET", "/x/ping", nil)
	if err != nil {
		t.Error(err)
	}
	if recorder.Code != 200 {
		t.Error("/x/ping status code:", recorder.Code)
	}
}
