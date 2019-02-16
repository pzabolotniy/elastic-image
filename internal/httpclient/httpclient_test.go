package httpclient

import (
	"errors"
	"fmt"
	"github.com/pzabolotniy/elastic-image/internal/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRestClientHandlerEnv_Get_OK(t *testing.T) {
	testName := t.Name()

	stubResponse := `{"ok":true}`
	httpTestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestMethod := r.Method

		expectedMethod := "GET"
		assert.Equalf(t, expectedMethod, requestMethod, "%s - http method ok", testName)

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(stubResponse))
	}))

	conf := config.GetConfig()
	logger := conf.APILogger
	timeout := conf.Timeout

	restClient := NewHTTPClient(logger, timeout)

	response, err := restClient.Get(httpTestServer.URL + "/ok")
	assert.NoErrorf(t, err, "%s - no error, ok", testName)

	expectedStatusCode := http.StatusOK
	gotStatusCode := response.StatusCode()
	assert.Equalf(t, expectedStatusCode, gotStatusCode, "%s - status code ok", testName)

	expectedBody := []byte(stubResponse)
	gotBody := response.RawBody()
	assert.Equalf(t, expectedBody, gotBody, "%s - body ok", testName)
}

func TestRestClientHandlerEnv_Get_NotFound(t *testing.T) {
	testName := t.Name()

	stubResponse := `not found`
	httpTestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestMethod := r.Method

		expectedMethod := "GET"
		assert.Equalf(t, expectedMethod, requestMethod, "%s - http method ok", testName)

		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "plain/text")
		w.Write([]byte(stubResponse))
	}))

	conf := config.GetConfig()
	logger := conf.APILogger
	timeout := conf.Timeout

	restClient := NewHTTPClient(logger, timeout)

	response, err := restClient.Get(httpTestServer.URL + "/not_found")
	expectedErr := errors.New(fmt.Sprintf("%d %s", http.StatusNotFound, http.StatusText(http.StatusNotFound)))
	assert.Equalf(t, expectedErr, err, "%s - no error, ok", testName)

	expectedStatusCode := http.StatusNotFound
	gotStatusCode := response.StatusCode()
	assert.Equalf(t, expectedStatusCode, gotStatusCode, "%s - status code ok", testName)

	expectedBody := []byte(stubResponse)
	gotBody := response.RawBody()
	assert.Equalf(t, expectedBody, gotBody, "%s - body ok", testName)
}

func TestRestClientHandlerEnv_Get_InternalServerError(t *testing.T) {
	testName := t.Name()

	stubResponse := `internal server error`
	httpTestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestMethod := r.Method

		expectedMethod := "GET"
		assert.Equalf(t, expectedMethod, requestMethod, "%s - http method ok", testName)

		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "plain/text")
		w.Write([]byte(stubResponse))
	}))

	conf := config.GetConfig()
	logger := conf.APILogger
	timeout := conf.Timeout

	restClient := NewHTTPClient(logger, timeout)

	response, err := restClient.Get(httpTestServer.URL + "/internal_server_error")
	expectedErr := errors.New(fmt.Sprintf("%d %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
	assert.Equalf(t, expectedErr, err, "%s - no error, ok", testName)

	expectedStatusCode := http.StatusInternalServerError
	gotStatusCode := response.StatusCode()
	assert.Equalf(t, expectedStatusCode, gotStatusCode, "%s - status code ok", testName)

	expectedBody := []byte(stubResponse)
	gotBody := response.RawBody()
	assert.Equalf(t, expectedBody, gotBody, "%s - body ok", testName)
}
