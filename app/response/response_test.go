package response

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestJSON(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	JSON(c, gin.H{"k": "v"})
	if c.Writer.Status() != 200 {
		t.Error("http status code error")
	}
	j := w.Body.Bytes()
	r := Response{}
	err := json.Unmarshal(j, &r)
	if err != nil {
		t.Error(err)
	}
	if r.Data.(map[string]interface{})["k"].(string) != "v" {
		t.Error("json data error")
	}
}

func TestJSON400(t *testing.T) {
	gin.SetMode(gin.DebugMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	ErrJSON400(c, "extra msg")
	if c.Writer.Status() != 400 {
		t.Error("http status code error")
	}
	j := w.Body.Bytes()
	r := Response{}
	err := json.Unmarshal(j, &r)
	if err != nil {
		t.Error(err)
	}
}

func TestJSON404(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	ErrJSON404(c, "extra msg")
	if c.Writer.Status() != 404 {
		t.Error("http status code error")
	}
	j := w.Body.Bytes()
	r := Response{}
	err := json.Unmarshal(j, &r)
	if err != nil {
		t.Error(err)
	}
}

func TestJSON500(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	ErrJSON500(c, "extra msg")
	if c.Writer.Status() != 500 {
		t.Error("http status code error")
	}
	j := w.Body.Bytes()
	r := Response{}
	err := json.Unmarshal(j, &r)
	if err != nil {
		t.Error(err)
	}
}

func TestRespond(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	Respond(c, 200, gin.H{"k": "v"}, RCSuccess)
	if c.Writer.Status() != 200 {
		t.Error("http status code error")
	}
	j := w.Body.Bytes()
	r := Response{}
	err := json.Unmarshal(j, &r)
	if err != nil {
		t.Error(err)
	}
	if r.Data.(map[string]interface{})["k"].(string) != "v" {
		t.Error("json data error")
	}

}
