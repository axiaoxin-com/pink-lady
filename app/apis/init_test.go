package apis

import (
	"testing"

	"github.com/axiaoxin/gin-skeleton/app/apis/router"
	"github.com/axiaoxin/gin-skeleton/app/utils"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestPing(t *testing.T) {
	r := router.SetupRouter("test", "", false)
	RegisterRoutes(r)
	w := utils.TestingGETRequest(r, "/x/ping")
	body := jsoniter.Get(w.Body.Bytes())
	version := body.Get("data", "version").ToString()
	if version != VERSION {
		t.Error("version error")
	}
}
