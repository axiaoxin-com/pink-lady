package apis

import (
	"testing"

	"pink-lady/app/router"
	"pink-lady/app/utils"

	jsoniter "github.com/json-iterator/go"
)

func TestPing(t *testing.T) {
	r := router.SetupRouter("../", "config")
	RegisterRoutes(r)
	w := utils.PerformRequest(r, "GET", "/x/ping", nil)
	body := jsoniter.Get(w.Body.Bytes())
	version := body.Get("data", "version").ToString()
	if version != VERSION {
		t.Error("version error", body.ToString())
	}
}
