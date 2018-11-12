package main

import (
	"os"
	"testing"

	"github.com/axiaoxin/gin-skeleton/app/utils"
	"github.com/gin-gonic/gin"
)

func setup() {
	gin.SetMode(gin.TestMode)
}

func teardown() {

}

func TestSetupAPP(t *testing.T) {
	app := SetupAPP()
	w := utils.PerformTestingRequest(app, "GET", "/x/ping")
	if w.Result().StatusCode != 200 {
		t.Error("SetupAPP fail")
	}
}

func TestMain(m *testing.M) {
	setup()
	exitCode := m.Run()
	teardown()
	os.Exit(exitCode)
}
