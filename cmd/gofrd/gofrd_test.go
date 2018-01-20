package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func test_handler_func(handler http.HandlerFunc, method, target string, body_reader io.Reader) (response *http.Response, body []byte, err error) {
	r := httptest.NewRequest(method, target, body_reader)
	w := httptest.NewRecorder()

	handler(w, r)
	response = w.Result()
	body, err = ioutil.ReadAll(response.Body)

	return
}

func TestVersion(t *testing.T) {
	assert.Regexp(t, `\d+\.\d+\.\d+`, Version)
}

func Test_route_status(t *testing.T) {
	response, body, err := test_handler_func(route_status, "GET", "http:///status", nil)

	assert.NoError(t, err)
	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, "OK", string(body))
}

func Test_route_png_query(t *testing.T) {
	target := "http:///png?i=100&w=100&h=100&e=4&m=%23444444&c=mono&r=mandelbrot&s=2&p=2&rmin=-2.6250000000000004&rmax=1.5750000000000002&imin=-2.1&imax=2.1&render_id=1fde6438-5aee-4c70-8780-cf321c638d8f"
	response, body, err := test_handler_func(route_png, "GET", target, nil)

	assert.NoError(t, err)
	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, []byte{0x89, 0x50, 0x4e, 0x47}, body[0:4])
}

func Test_route_png_no_query(t *testing.T) {
	response, body, err := test_handler_func(route_png, "GET", "http:///status", nil)

	assert.NoError(t, err)
	assert.Equal(t, response.StatusCode, http.StatusUnprocessableEntity)
	assert.Regexp(t, `(?i)invalid`, string(body))
}
