package main

import (
	"testing"

	"github.com/axiaoxin/gin-skeleton/app/utils"
)

func TestSetupAPP(t *testing.T) {
	app := SetupAPP()
	w := utils.PerformTestingRequest(app, "GET", "/x/ping")
	if w.Result().StatusCode != 200 {
		t.Error("SetupAPP fail")
	}
}
