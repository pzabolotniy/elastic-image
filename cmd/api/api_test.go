package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pzabolotniy/elastic-image/internal/config"
	"github.com/pzabolotniy/elastic-image/internal/httpclient"
	httpMocks "github.com/pzabolotniy/elastic-image/internal/httpclient/mocks"
	"github.com/pzabolotniy/elastic-image/internal/image/resize"
	resizeMocks "github.com/pzabolotniy/elastic-image/internal/image/resize/mocks"
	"github.com/pzabolotniy/elastic-image/internal/logging"
	"github.com/pzabolotniy/monkey"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestPostImageResize(t *testing.T) {
	router := gin.Default()
	testConfig := config.GetConfig()
	var nilError error

	setupRouter(router, testConfig)

	testRequest := struct {
		method string
		uri    string
	}{method: "POST", uri: "/api/v1/images/resize"}
	testName := t.Name()

	url := "http://any-host.local/image.jpeg"
	width := 1024
	height := 800

	mockedImage := "i.am.image"
	mockedBrowser := &httpMocks.Browser{}
	mockedResponse := &httpMocks.Responser{}
	monkey.Patch(httpclient.NewHTTPClient, func(logger logging.Logger, timeout time.Duration) httpclient.Browser {
		return mockedBrowser
	})

	imageBytes := []byte(mockedImage)
	mockedBrowser.On("Get", url).Return(mockedResponse, nilError)
	mockedResponse.On("RawBody").Return(imageBytes)
	imageReader := bytes.NewReader(imageBytes)

	mockedResizer := &resizeMocks.Resizer{}
	monkey.Patch(resize.NewResizer, func(logger logging.Logger) resize.Resizer {
		return mockedResizer
	})
	mockedResizer.On("Resize", imageReader, uint(width), uint(height)).Return(imageBytes, nilError)

	inputBody := fmt.Sprintf(`{"url":"%s","width":%d,"heigth":%d}`, url, width, height)
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

	mockedBrowser.AssertExpectations(t)
	mockedResponse.AssertExpectations(t)
	mockedResizer.AssertExpectations(t)
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
