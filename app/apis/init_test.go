package apis

import (
	"net/http/httptest"
	"testing"

	"github.com/axiaoxin/gin-skeleton/app/services"
	"github.com/axiaoxin/gin-skeleton/app/utils"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestPing(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	Ping(c)
	body := jsoniter.Get(w.Body.Bytes())
	version := body.Get("data", "version").ToString()
	if version != services.VERSION {
		t.Error("version error")
	}
}

func TestSetupRouter(t *testing.T) {
	app := SetupRouter("test", "", false)
	w := utils.PerformTestingRequest(app, "GET", "/x/ping")
	if w.Result().StatusCode != 200 {
		t.Error("Setup router fail")
	}
}
