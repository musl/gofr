package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testHandlerFunc(handler http.HandlerFunc, method, target string, bodyReader io.Reader) (response *http.Response, body []byte, err error) {
	r := httptest.NewRequest(method, target, bodyReader)
	w := httptest.NewRecorder()

	handler(w, r)
	response = w.Result()
	body, err = ioutil.ReadAll(response.Body)

	return
}

func TestVersion(t *testing.T) {
	assert.Regexp(t, `\d+\.\d+\.\d+`, Version)
}

func TestRouteStatus(t *testing.T) {
	response, body, err := testHandlerFunc(routeStatus, "GET", "http:///status", nil)

	assert.NoError(t, err)
	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, "OK", string(body))
}

func TestRoutePNGQuery(t *testing.T) {
	target := "http:///png?i=100&w=100&h=100&e=4&m=%23444444&c=mono&r=mandelbrot&s=2&p=2&rmin=-2.6250000000000004&rmax=1.5750000000000002&imin=-2.1&imax=2.1&render-id=1fde6438-5aee-4c70-8780-cf321c638d8f"
	response, body, err := testHandlerFunc(routePNG, "GET", target, nil)

	assert.NoError(t, err)
	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, []byte{0x89, 0x50, 0x4e, 0x47}, body[0:4])
}

func TestRoutePNGNoQuery(t *testing.T) {
	response, body, err := testHandlerFunc(routePNG, "GET", "http:///status", nil)

	assert.NoError(t, err)
	assert.Equal(t, response.StatusCode, http.StatusUnprocessableEntity)
	assert.Regexp(t, `(?i)invalid`, string(body))
}
