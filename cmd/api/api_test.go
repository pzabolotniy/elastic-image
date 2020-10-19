package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pzabolotniy/elastic-image/internal/config"
	"github.com/pzabolotniy/elastic-image/internal/image/resize"
	"github.com/pzabolotniy/elastic-image/internal/logging"
	"github.com/pzabolotniy/monkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type EnvTestSuite struct {
	suite.Suite
	router *gin.Engine
	conf   *config.AppConfig
}

func TestPostImageResize(t *testing.T) {
	envSuite := new(EnvTestSuite)
	suite.Run(t, envSuite)
}

func (s *EnvTestSuite) SetupSuite() {
	router := gin.New()
	testConfig := &config.AppConfig{
		ServerConfig: &config.ServerConfig{
			Bind: "",
		},
		ImageConfig: &config.ImageConfig{
			BrowserCacheTTL: 3600,
			FetchTimeout:    10 * time.Second,
		},
	}
	testLogger := logging.GetLogger()
	setupRouter(router, testConfig, testLogger)
	s.router = router
	s.conf = testConfig
}

func (s *EnvTestSuite) TestPostImageResize_OK() {
	t := s.T()
	router := s.router
	testRequest := struct {
		method string
		uri    string
	}{method: "POST", uri: "/api/v1/images/resize"}

	expectedWidth := 1024
	expectedHeight := 800

	mockedImage := "i.am.image"
	testURI := "/image.jpeg"

	httpTestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestMethod := r.Method
		requestPath := r.URL.String()

		expectedMethod := "GET"
		assert.Equal(t, expectedMethod, requestMethod, "http method ok")
		expectedPath := testURI
		assert.Equal(t, expectedPath, requestPath, "request path ok")

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "image/jpeg")
		_, err := w.Write([]byte(mockedImage))
		if err != nil {
			t.Fatal(err)
		}
	}))
	url := httpTestServer.URL + testURI

	monkey.Patch(resize.Resize, func(ctx context.Context, srcImage io.Reader, width, height uint) ([]byte, error) {
		assert.Equal(t, expectedWidth, int(width), "resize.Resize width ok")
		assert.Equal(t, expectedHeight, int(height), "resize.Resize height ok")

		expectedImage := bytes.NewReader([]byte(mockedImage))
		assert.Equal(t, expectedImage, srcImage, "resize.Resize image ok")
		return []byte(mockedImage), nil
	})

	inputBody := fmt.Sprintf(`{"url":"%s","width":%d,"heigth":%d}`, url, expectedWidth, expectedHeight)
	reader := strings.NewReader(inputBody)
	headers := make(map[string][]string)
	apiResponse := makeTestRequest(router, testRequest.method, testRequest.uri, reader, headers)

	expectedCode := http.StatusOK
	assert.Equal(t, expectedCode, apiResponse.Code, "http code ok")
	gotBody, err := ioutil.ReadAll(apiResponse.Body)
	if err != nil {
		t.Fatal(err)
	}
	expectedBody := []byte(mockedImage)
	assert.Equal(t, expectedBody, gotBody, "response ok")

	apiHeaders := apiResponse.Header()
	contentLength := len(mockedImage)
	expectedHeaders := http.Header{
		"Content-Type":   []string{"image/jpeg"},
		"Content-Length": []string{fmt.Sprintf("%d", contentLength)},
		"Cache-Control":  []string{fmt.Sprintf("max-age=%d, public", s.conf.ImageConfig.BrowserCacheTTL)},
		"X-Request-Id":   []string{apiHeaders.Get("X-Request-Id")}, // will not assert random uuid
	}
	assert.Equal(t, expectedHeaders, apiHeaders, "response headers ok")
}

func makeTestRequest(r http.Handler, method, path string, body io.Reader, headers map[string][]string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	httpHeaders := make(http.Header)
	for k, v := range headers {
		httpHeaders[k] = v
	}
	req.Header = httpHeaders
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
