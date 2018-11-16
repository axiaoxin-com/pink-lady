package utils

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/alicebob/miniredis"
	"github.com/sirupsen/logrus"
)

// TestingGETRequest perform a GET request with the handler for testing
func TestingGETRequest(r http.Handler, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// TestingPOSTRequest perform a POST request with the handler for testing
func TestingPOSTRequest(r http.Handler, path string, jsonStr string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", path, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// MockRedis provide a mock redis server for testing
func MockRedis() (*miniredis.Miniredis, error) {
	s, err := miniredis.Run()
	if err != nil {
		logrus.Error(err)
	}
	return s, err
}
