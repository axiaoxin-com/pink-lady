package response

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/axiaoxin/gin-skeleton/app/services/retcode"
	"github.com/gin-gonic/gin"
)

func init() {
	// 避免打印那些个debuglog
	gin.SetMode("release")
}

func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	JSON(c, retcode.Success, gin.H{"k": "v"})
	if c.Writer.Status() != 200 {
		t.Error("http status code error")
	}
	j := w.Body.Bytes()
	r := Response{}
	err := json.Unmarshal(j, &r)
	if err != nil {
		t.Error(err)
	}
	if r.Code != 0 {
		t.Error("json code error")
	}
	if r.Data.(map[string]interface{})["k"].(string) != "v" {
		t.Error("json data error")
	}
}

func TestJSON400(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	JSON400(c, retcode.InvalidParams, gin.H{"k": "v"})
	if c.Writer.Status() != 400 {
		t.Error("http status code error")
	}
	j := w.Body.Bytes()
	r := Response{}
	err := json.Unmarshal(j, &r)
	if err != nil {
		t.Error(err)
	}
	if r.Code != 3 {
		t.Error("json code error")
	}
	if r.Data.(map[string]interface{})["k"].(string) != "v" {
		t.Error("json data error")
	}
}

func TestJSON404(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	JSON404(c, retcode.RouteNotFound, gin.H{"k": "v"})
	if c.Writer.Status() != 404 {
		t.Error("http status code error")
	}
	j := w.Body.Bytes()
	r := Response{}
	err := json.Unmarshal(j, &r)
	if err != nil {
		t.Error(err)
	}
	if r.Code != 4 {
		t.Error("json code error")
	}
	if r.Data.(map[string]interface{})["k"].(string) != "v" {
		t.Error("json data error")
	}
}

func TestJSON500(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	JSON500(c, retcode.InternalError, gin.H{"k": "v"})
	if c.Writer.Status() != 500 {
		t.Error("http status code error")
	}
	j := w.Body.Bytes()
	r := Response{}
	err := json.Unmarshal(j, &r)
	if err != nil {
		t.Error(err)
	}
	if r.Code != 5 {
		t.Error("json code error")
	}
	if r.Data.(map[string]interface{})["k"].(string) != "v" {
		t.Error("json data error")
	}
}
