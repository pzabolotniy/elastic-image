package httpclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pzabolotniy/elastic-image/internal/logging"
	"github.com/stretchr/testify/assert"
)

func TestRestClientHandlerEnv_Get_OK(t *testing.T) {
	stubResponse := `{"ok":true}`
	httpTestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestMethod := r.Method

		expectedMethod := "GET"
		assert.Equal(t, expectedMethod, requestMethod, "http method ok")

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(stubResponse))
	}))

	logger := logging.GetLogger()
	fetchTimeout := 10 * time.Second
	restClient := NewHTTPClient(logger, fetchTimeout)

	response, err := restClient.Get(httpTestServer.URL + "/ok")
	assert.NoError(t, err, "no error, ok")

	expectedStatusCode := http.StatusOK
	gotStatusCode := response.StatusCode()
	assert.Equal(t, expectedStatusCode, gotStatusCode, "status code ok")

	expectedBody := []byte(stubResponse)
	gotBody := response.RawBody()
	assert.Equal(t, expectedBody, gotBody, "%s - body ok")
}

func TestRestClientHandlerEnv_Get_NotFound(t *testing.T) {
	stubResponse := `not found`
	httpTestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestMethod := r.Method

		expectedMethod := "GET"
		assert.Equalf(t, expectedMethod, requestMethod, "http method ok")

		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "plain/text")
		w.Write([]byte(stubResponse))
	}))

	logger := logging.GetLogger()
	fetchTimeout := 10 * time.Second
	restClient := NewHTTPClient(logger, fetchTimeout)

	response, err := restClient.Get(httpTestServer.URL + "/not_found")
	expectedErr := fmt.Errorf("%d %s", http.StatusNotFound, http.StatusText(http.StatusNotFound))
	assert.Equal(t, expectedErr, err, "no error, ok")

	expectedStatusCode := http.StatusNotFound
	gotStatusCode := response.StatusCode()
	assert.Equal(t, expectedStatusCode, gotStatusCode, "status code ok")

	expectedBody := []byte(stubResponse)
	gotBody := response.RawBody()
	assert.Equal(t, expectedBody, gotBody, "%s - body ok")
}

func TestRestClientHandlerEnv_Get_InternalServerError(t *testing.T) {
	stubResponse := `internal server error`
	httpTestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestMethod := r.Method

		expectedMethod := "GET"
		assert.Equal(t, expectedMethod, requestMethod, "http method ok")

		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "plain/text")
		w.Write([]byte(stubResponse))
	}))

	fetchTimeout := 10 * time.Second
	logger := logging.GetLogger()
	restClient := NewHTTPClient(logger, fetchTimeout)

	response, err := restClient.Get(httpTestServer.URL + "/internal_server_error")
	expectedErr := fmt.Errorf("%d %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	assert.Equal(t, expectedErr, err, "%s - no error, ok")

	expectedStatusCode := http.StatusInternalServerError
	gotStatusCode := response.StatusCode()
	assert.Equal(t, expectedStatusCode, gotStatusCode, "%s - status code ok")

	expectedBody := []byte(stubResponse)
	gotBody := response.RawBody()
	assert.Equal(t, expectedBody, gotBody, "%s - body ok")
}
