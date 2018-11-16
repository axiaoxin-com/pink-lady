package apis

import (
	"testing"

	"pink-lady/app/utils"

	"github.com/gin-gonic/gin"
)

func TestRegisterRoutes(t *testing.T) {
	r := gin.New()
	RegisterRoutes(r)
	w := utils.TestingGETRequest(r, "/x/ping")
	if w.Result().StatusCode != 200 {
		t.Error("register routes no /x/ping")
	}
}
