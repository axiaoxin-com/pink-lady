package router

import (
	"testing"

	"pink-lady/app/utils"

	"github.com/gin-gonic/gin"
)

func TestSetupRouter(t *testing.T) {
	r := SetupRouter("..", "config")
	r.GET("/", func(c *gin.Context) {
		return
	})
	respRecorder := utils.PerformRequest(r, "GET", "/", nil)
	if respRecorder.Result().StatusCode != 200 {
		t.Fatal("Setup router fail")
	}
}
