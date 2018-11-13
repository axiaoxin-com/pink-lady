package router

import (
	"testing"

	"github.com/axiaoxin/gin-skeleton/app/utils"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestSetupRouter(t *testing.T) {
	r := SetupRouter("test", "", false)
	r.GET("/xyz", func(c *gin.Context) {})
	w := utils.PerformTestingRequest(r, "GET", "/xyz")
	if w.Result().StatusCode != 200 {
		t.Error("Setup router fail")
	}
}
