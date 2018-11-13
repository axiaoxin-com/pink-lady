package utils

import (
	"bytes"
	"net/http"
	"net/http/httptest"
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
