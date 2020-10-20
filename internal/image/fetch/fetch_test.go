package fetch

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pzabolotniy/elastic-image/internal/logging"
	"github.com/stretchr/testify/assert"
)

func TestNewFetchParams(t *testing.T) {
	timeout := 10 * time.Second
	url := "http://any-host.local/image.jpeg"

	params := NewFetchParams(timeout, url)
	expected := &Params{
		Timeout: timeout,
		URL:     url,
	}
	assert.Equal(t, expected, params, "result ok")
}

func TestGetImage(t *testing.T) {
	ctx := context.Background()
	testLogger := logging.GetLogger()
	ctx = logging.WithContext(ctx, testLogger)
	timeout := 10 * time.Second
	testURI := "/image.jpeg"
	sharedDownloads := make(map[string]*DownloadState)

	mockedImage := []byte("i.am.image")
	httpTestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestMethod := r.Method
		requestPath := r.URL.String()

		expectedMethod := "GET"
		assert.Equal(t, expectedMethod, requestMethod, "http method ok")
		expectedPath := testURI
		assert.Equal(t, expectedPath, requestPath, "request path ok")

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "image/jpeg")
		_, err := w.Write(mockedImage)
		if err != nil {
			t.Fatal(err)
		}
	}))
	url := httpTestServer.URL + testURI

	params := NewFetchParams(timeout, url)
	response, err := GetImage(ctx, sharedDownloads, params)

	assert.NoError(t, err, "no errors, ok")
	gotImage, err := ioutil.ReadAll(response)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, mockedImage, gotImage, "image bytes ok")
}

func TestGetImage_SharedDownload(t *testing.T) {
	ctx := context.Background()
	testLogger := logging.GetLogger()
	ctx = logging.WithContext(ctx, testLogger)
	timeout := 10 * time.Second
	testURI := "/image.jpeg"
	sharedDownloads := make(map[string]*DownloadState)

	mockedImage := []byte("i.am.image")
	fetchCounter := 0
	expectedFetchCounter := 1
	httpTestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fetchCounter++
		requestMethod := r.Method
		requestPath := r.URL.String()

		expectedMethod := "GET"
		assert.Equal(t, expectedMethod, requestMethod, "http method ok")
		expectedPath := testURI
		assert.Equal(t, expectedPath, requestPath, "request path ok")

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "image/jpeg")
		_, err := w.Write(mockedImage)
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(100 * time.Millisecond)
	}))
	url := httpTestServer.URL + testURI

	params := NewFetchParams(timeout, url)
	var (
		resp1, resp2 io.Reader
		err1, err2   error
	)

	go func() {
		resp1, err1 = GetImage(ctx, sharedDownloads, params)
	}()
	time.Sleep(50 * time.Millisecond)
	resp2, err2 = GetImage(ctx, sharedDownloads, params)

	assert.NoError(t, err1, "no errors 1, ok")
	gotImage1, err := ioutil.ReadAll(resp1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, mockedImage, gotImage1, "image 1 bytes ok")

	assert.NoError(t, err2, "no errors 2, ok")
	gotImage2, err := ioutil.ReadAll(resp2)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, mockedImage, gotImage2, "image 2 bytes ok")

	assert.Equal(t, expectedFetchCounter, fetchCounter, "http fetch requests count ok")
}

func TestGetImage_SharedDownload_FetchFailed(t *testing.T) {
	ctx := context.Background()
	testLogger := logging.GetLogger()
	ctx = logging.WithContext(ctx, testLogger)
	timeout := 10 * time.Second
	testURI := "/image.jpeg"
	sharedDownloads := make(map[string]*DownloadState)

	fetchCounter := 0
	expectedFetchCounter := 1
	httpTestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fetchCounter++
		requestMethod := r.Method
		requestPath := r.URL.String()

		expectedMethod := "GET"
		assert.Equal(t, expectedMethod, requestMethod, "http method ok")
		expectedPath := testURI
		assert.Equal(t, expectedPath, requestPath, "request path ok")

		w.WriteHeader(http.StatusInternalServerError)
		time.Sleep(100 * time.Millisecond)
	}))
	url := httpTestServer.URL + testURI

	params := NewFetchParams(timeout, url)
	var (
		resp1, resp2 io.Reader
		err1, err2   error
	)

	go func() {
		resp1, err1 = GetImage(ctx, sharedDownloads, params)
	}()
	time.Sleep(50 * time.Millisecond)
	resp2, err2 = GetImage(ctx, sharedDownloads, params)

	expectedError := errors.New(fmt.Sprintf("%d %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
	assert.Equal(t, expectedError, err1, "error 1 ok")
	assert.Equal(t, nil, resp1, "image 1 nil, ok")

	assert.Equal(t, expectedError, err2, "error 2 ok")
	assert.Equal(t, nil, resp2, "image 2 nil, ok")

	assert.Equal(t, expectedFetchCounter, fetchCounter, "http fetch requests count ok")
}
