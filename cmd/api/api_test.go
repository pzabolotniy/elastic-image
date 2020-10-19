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
	"github.com/pzabolotniy/elastic-image/internal/httpclient"
	httpMocks "github.com/pzabolotniy/elastic-image/internal/httpclient/mocks"
	"github.com/pzabolotniy/elastic-image/internal/image/resize"
	"github.com/pzabolotniy/elastic-image/internal/logging"
	"github.com/pzabolotniy/monkey"
	"github.com/stretchr/testify/assert"
)

func TestPostImageResize(t *testing.T) {
	router := gin.Default()
	testConfig := &config.AppConfig{
		ServerConfig: &config.ServerConfig{
			Bind: "",
		},
		ImageConfig: &config.ImageConfig{
			CacheTTL:     3600 * time.Second,
			FetchTimeout: 10 * time.Second,
		},
	}
	testLogger := logging.GetLogger()
	var nilError error

	setupRouter(router, testConfig, testLogger)

	testRequest := struct {
		method string
		uri    string
	}{method: "POST", uri: "/api/v1/images/resize"}
	testName := t.Name()

	url := "http://any-host.local/image.jpeg"
	expectedWidth := 1024
	expectedHeight := 800

	mockedImage := "i.am.image"
	mockedBrowser := &httpMocks.Browser{}
	mockedResponse := &httpMocks.Responser{}
	monkey.Patch(httpclient.NewHTTPClient, func(logger logging.Logger, timeout time.Duration) httpclient.Browser {
		return mockedBrowser
	})

	imageBytes := []byte(mockedImage)
	mockedBrowser.On("Get", url).Return(mockedResponse, nilError)
	mockedResponse.On("RawBody").Return(imageBytes)
	expectedImage := bytes.NewReader(imageBytes)

	//mockedResizer := &resizeMocks.Resizer{}
	monkey.Patch(resize.Resize, func(ctx context.Context, srcImage io.Reader, width, height uint) ([]byte, error) {
		assert.Equal(t, expectedWidth, int(width), "resize.Resize width ok")
		assert.Equal(t, expectedHeight, int(height), "resize.Resize height ok")
		assert.Equal(t, expectedImage, srcImage, "resize.Resize image ok")
		return imageBytes, nil
	})

	inputBody := fmt.Sprintf(`{"url":"%s","width":%d,"heigth":%d}`, url, expectedWidth, expectedHeight)
	reader := strings.NewReader(inputBody)
	headers := make(map[string][]string)
	apiResponse := makeTestRequest(router, testRequest.method, testRequest.uri, reader, headers)

	expectedCode := http.StatusOK
	assert.Equalf(t, expectedCode, apiResponse.Code, "%s - http code ok", testName)
	gotBodybytes, err := ioutil.ReadAll(apiResponse.Body)
	if err != nil {
		panic(err)
	}
	expectedBodyBytes := imageBytes
	assert.Equalf(t, expectedBodyBytes, gotBodybytes, "%s - response ok", testName)
	apiHeaders := apiResponse.Header()
	contentLength := len(mockedImage)
	cacheTTL := testConfig.ImageConfig.CacheTTL
	expectedHeaders := http.Header{
		"Content-Type":   []string{"image/jpeg"},
		"Content-Length": []string{fmt.Sprintf("%d", contentLength)},
		"Cache-Control":  []string{fmt.Sprintf("max-age=%d, public", cacheTTL)},
	}
	assert.Equalf(t, expectedHeaders, apiHeaders, "%s - response headers ok", testName)

	mockedBrowser.AssertExpectations(t)
	mockedResponse.AssertExpectations(t)
	//mockedResizer.AssertExpectations(t)
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
