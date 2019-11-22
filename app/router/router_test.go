package router

import (
	"testing"

	"github.com/axiaoxin/pink-lady/app/logging"
	"github.com/axiaoxin/pink-lady/app/utils"

	"github.com/gin-gonic/gin"
)

func TestSetupRouter(t *testing.T) {
	logging.InitLogger()
	r := SetupRouter()
	r.GET("/xyz", func(c *gin.Context) {})
	w := utils.TestingGETRequest(r, "/xyz")
	if w.Result().StatusCode != 200 {
		t.Error("Setup router fail")
	}
}
