package router

import (
	"testing"

	"pink-lady/app/utils"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestSetupRouter(t *testing.T) {
	r := SetupRouter("test", "", false)
	r.GET("/xyz", func(c *gin.Context) {})
	w := utils.TestingGETRequest(r, "/xyz")
	if w.Result().StatusCode != 200 {
		t.Error("Setup router fail")
	}
	if gin.Mode() != gin.TestMode {
		t.Error("Set mode fail")
	}
	r = SetupRouter("debug", "", false)
	if gin.Mode() != gin.DebugMode {
		t.Error("Set mode fail")
	}
	r = SetupRouter("release", "no-sentry", false)
	if gin.Mode() != gin.ReleaseMode {
		t.Error("Set mode fail")
	}
}
