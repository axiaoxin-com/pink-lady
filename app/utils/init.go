// package utils save the package of third party package tools and the general tool code written by yourself
// write your general tool in the package
package utils

import (
	"net/http"
	"net/http/httptest"
)

// init function here for utils package
func init() {

}

// PerformTestingRequest perform a request with the handler for testing
func PerformTestingRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
