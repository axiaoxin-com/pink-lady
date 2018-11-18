package middleware

import (
	"bytes"
	"encoding/json"
	"testing"

	"pink-lady/app/utils"

	"github.com/gin-gonic/gin"
)

type logStruct struct {
	Path      string `json:"path"`
	Method    string `json:"method"`
	ClientIP  string `json:"clientIP"`
	UserAgent string `json:"userAgent"`
	RequestID string `json:"requestID"`
	Status    int    `json:"status"`
	Size      int    `json:"size"`
	Latency   string `json:"latency"`
	Level     string `json:"level"`
	Msg       string `json:"msg"`
	Time      string `json:"time"`
}

func TestGinLogrus(t *testing.T) {
	buf := &bytes.Buffer{}
	ls := logStruct{}
	utils.InitLogrus(buf, "debug", "json")
	r := gin.New()
	r.Use(RequestID()) // must get request id
	r.Use(GinLogrus())
	r.GET("/", func(c *gin.Context) {})
	utils.TestingGETRequest(r, "/")
	err := json.Unmarshal(buf.Bytes(), &ls)
	if err != nil {
		t.Error(err)
	}
	if ls.Path != "/" || ls.Method != "GET" || ls.RequestID == "" || ls.Status != 200 {
		t.Error("log field error ", ls)
	}
}

func TestGinLogrus500(t *testing.T) {
	buf := &bytes.Buffer{}
	ls := logStruct{}
	utils.InitLogrus(buf, "debug", "json")
	r := gin.New()
	r.Use(RequestID()) // must get request id
	r.Use(GinLogrus())
	r.GET("/500", func(c *gin.Context) { c.AbortWithStatus(500) })
	utils.TestingGETRequest(r, "/500")
	err := json.Unmarshal(buf.Bytes(), &ls)
	if err != nil {
		t.Error(err)
	}
	if ls.Path != "/500" || ls.Method != "GET" || ls.RequestID == "" || ls.Status != 500 {
		t.Error("log field error ", ls)
	}
}
