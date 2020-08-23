package webserver

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/axiaoxin-com/goutils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func TestViperConfig(t *testing.T) {
	InitViperConfig("..", "config.default", "toml")
	defer viper.Reset()
	if !goutils.IsInitedViper() {
		t.Error("init viper failed")
	}
}

func TestNewGinEngine(t *testing.T) {
	viper.SetDefault("server.mode", "release")
	defer viper.Reset()
	if app := NewGinEngine(); app == nil {
		t.Error("app is nil")
	}
}

func TestRun(t *testing.T) {
	InitViperConfig("..", "config.default", "toml")
	defer viper.Reset()
	viper.Set("server.mode", "release")
	viper.Set("logging.level", "error")
	app := NewGinEngine()
	register := func(g http.Handler) {
		g.(*gin.Engine).GET("/", func(c *gin.Context) {
			c.JSON(200, 666)
		})
	}
	go Run(app, register)
	rsp, err := http.Get("http://localhost" + viper.GetString("server.addr"))
	if err != nil {
		t.Fatal("request running server error:", err)
	}
	defer rsp.Body.Close()
	if b, err := ioutil.ReadAll(rsp.Body); err != nil {
		t.Error("read running server response body error:", err)
	} else if string(b) != "666" {
		t.Error("running server response invalid:", string(b))
	}
}
