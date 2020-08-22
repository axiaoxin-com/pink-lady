package apis

import (
	"testing"

	"github.com/axiaoxin-com/goutils"
	"github.com/gin-gonic/gin"
)

func TestRegisterRoutes(t *testing.T) {
	r := gin.New()
	Register(r)
	recorder, err := goutils.RequestHTTPHandler(r, "GET", "/x/ping", nil)
	if err != nil {
		t.Error(err)
	}
	if recorder.Code != 200 {
		t.Error("/x/ping status code:", recorder.Code)
	}
}
