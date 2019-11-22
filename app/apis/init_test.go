package apis

import (
	"io"
	"os"
	"testing"

	"github.com/axiaoxin/pink-lady/app/logging"
	"github.com/axiaoxin/pink-lady/app/router"
	"github.com/axiaoxin/pink-lady/app/utils"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

func TestPing(t *testing.T) {
	// 配置文件默认加载当前目录，需要把配置文件移到这里
	confile, err := os.Open("../config.toml.example")
	if err != nil {
		t.Error(err)
	}
	defer confile.Close()
	tmpConfile, err := os.Create("./config.toml")
	if err != nil {
		t.Error(err)
	}
	defer tmpConfile.Close()
	io.Copy(tmpConfile, confile)
	// 清理测试用的配置文件
	defer func() { os.Remove("./config.toml") }()
	workdir, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	utils.InitViper(workdir, "config", "envPrefix")
	logging.InitLogger()

	gin.SetMode(gin.TestMode)
	r := router.SetupRouter()
	RegisterRoutes(r)
	w := utils.TestingGETRequest(r, "/x/ping")
	body := jsoniter.Get(w.Body.Bytes())
	version := body.Get("data", "version").ToString()
	if version != VERSION {
		t.Error("version error")
	}
}
