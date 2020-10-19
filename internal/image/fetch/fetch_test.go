package fetch

import (
	"context"
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
	response, err := GetImage(ctx, params)

	assert.NoError(t, err, "no errors, ok")
	gotImage, err := ioutil.ReadAll(response)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, mockedImage, gotImage, "image bytes ok")
}
