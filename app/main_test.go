package main

import (
	"testing"

	"github.com/axiaoxin/gin-skeleton/app/apis"
	"github.com/axiaoxin/gin-skeleton/app/utils"
)

func TestSetupAPP(t *testing.T) {
	app := apis.SetupRouter("test", "", false)
	w := utils.PerformTestingRequest(app, "GET", "/x/ping")
	if w.Result().StatusCode != 200 {
		t.Error("SetupAPP fail")
	}
}
