package response

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestJSON(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	responseWriter := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(responseWriter)
	JSON(c, gin.H{"k": "v"})
	if c.Writer.Status() != 200 {
		t.Fatal("http status code error")
	}
	j := responseWriter.Body.Bytes()
	r := Response{}
	err := json.Unmarshal(j, &r)
	if err != nil {
		t.Fatal(err)
	}
	if r.Data.(map[string]interface{})["k"].(string) != "v" {
		t.Fatal("json data error")
	}
}

func TestErrJSON(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	responseWriter := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(responseWriter)
	ErrJSON(c, CodeInvalidParam)
	if c.Writer.Status() != 200 {
		t.Fatal("http status code error")
	}
	j := responseWriter.Body.Bytes()
	r := Response{}
	err := json.Unmarshal(j, &r)
	if err != nil {
		t.Fatal(err)
	}
}

func TestJSON400(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	responseWriter := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(responseWriter)
	ErrJSON400(c, "extra msg")
	if c.Writer.Status() != 400 {
		t.Fatal("http status code error")
	}
	j := responseWriter.Body.Bytes()
	r := Response{}
	err := json.Unmarshal(j, &r)
	if err != nil {
		t.Fatal(err)
	}
}

func TestJSON404(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	responseWriter := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(responseWriter)
	ErrJSON404(c, "extra msg")
	if c.Writer.Status() != 404 {
		t.Fatal("http status code error")
	}
	j := responseWriter.Body.Bytes()
	r := Response{}
	err := json.Unmarshal(j, &r)
	if err != nil {
		t.Fatal(err)
	}
}

func TestJSON500(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	responseWriter := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(responseWriter)
	ErrJSON500(c, "extra msg")
	if c.Writer.Status() != 500 {
		t.Fatal("http status code error")
	}
	j := responseWriter.Body.Bytes()
	r := Response{}
	err := json.Unmarshal(j, &r)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRespond(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	responseWriter := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(responseWriter)
	Respond(c, 200, gin.H{"k": "v"}, CodeSuccess)
	if c.Writer.Status() != 200 {
		t.Fatal("http status code error")
	}
	j := responseWriter.Body.Bytes()
	r := Response{}
	err := json.Unmarshal(j, &r)
	if err != nil {
		t.Fatal(err)
	}
	if r.Data.(map[string]interface{})["k"].(string) != "v" {
		t.Fatal("json data error")
	}
}

func TestRespondWithExtraMsg(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	responseWriter := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(responseWriter)
	Respond(c, 200, gin.H{"k": "v"}, CodeSuccess, "xxx")
	if c.Writer.Status() != 200 {
		t.Fatal("http status code error")
	}
	j := responseWriter.Body.Bytes()
	r := Response{}
	err := json.Unmarshal(j, &r)
	if err != nil {
		t.Fatal(err)
	}
	if r.Data.(map[string]interface{})["k"].(string) != "v" {
		t.Fatal("json data error")
	}
	if !strings.Contains(r.Msg, "xxx") {
		t.Fatal("extraMsgs error", r)
	}
}

func TestRespondWithError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	responseWriter := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(responseWriter)
	Respond(c, 200, gin.H{"k": "v"}, errors.New("errxx"), "xxx")
	if c.Writer.Status() != 200 {
		t.Fatal("http status code error")
	}
	j := responseWriter.Body.Bytes()
	r := Response{}
	err := json.Unmarshal(j, &r)
	if err != nil {
		t.Fatal(err)
	}
	if r.Data.(map[string]interface{})["k"].(string) != "v" {
		t.Fatal("json data error")
	}
	if !strings.Contains(r.Msg, "errxx") {
		t.Fatal("extraMsgs error", r)
	}
}
