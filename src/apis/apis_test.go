package apis

import (
	"testing"

	"pink-lady/app/utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func TestRegisterRoutes(t *testing.T) {
	r := gin.New()
	RegisterRoutes(r)
	w := utils.PerformRequest(r, "GET", "/x/ping", nil)
	if w.Result().StatusCode != 200 {
		t.Fatal("register routes no /x/ping")
	}
}

func TestRedirect(t *testing.T) {
	r := gin.New()
	RegisterRoutes(r)
	viper.SetDefault("apidocs.rootRedirect", true)
	w := utils.PerformRequest(r, "GET", "/", nil)
	if w.Result().StatusCode != 301 {
		t.Fatal("root redirect fail")
	}
}
