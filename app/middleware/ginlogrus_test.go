package middleware

import (
	"bytes"
	"encoding/json"
	"testing"

	"pink-lady/app/utils"

	"github.com/gin-gonic/gin"
)

func TestGinLogrus(t *testing.T) {
	type logStruct struct {
		Path      string `json:"path"`
		Method    string `json:"method"`
		ClientIP  string `json:"clientIP"`
		UserAgent string `json:"userAgent"`
		RequestID string `json:"requestID"`
		Status    int    `json:"status"`
		Size      int    `json:"size"`
		Latency   string `json:"latency"`
	}
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
