package utils

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/alicebob/miniredis"
)

// PerformRequest perform a request with the handler for testing
func PerformRequest(app http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	method = strings.ToUpper(method)
	req, _ := http.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	respRecorder := httptest.NewRecorder()
	app.ServeHTTP(respRecorder, req)
	return respRecorder
}

// MockRedis provide a mock redis server for testing
func MockRedis() (*miniredis.Miniredis, error) {
	s, err := miniredis.Run()
	if err != nil {
		log.Println("[Error] ", err)
	}
	return s, err
}

// MockHTTPServer provide a mock http server for testing
// Usage: https://gist.github.com/axiaoxin/aa8014738c6a02ce4e66eb01168d24fe
func MockHTTPServer(f http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(f))
}

// Assert fails the test if the condition is false.
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// OK fails the test if an err is not nil.
func OK(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// Equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
