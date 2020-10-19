package fetch

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"github.com/pzabolotniy/elastic-image/internal/config"
	"github.com/pzabolotniy/elastic-image/internal/httpclient"
	"github.com/pzabolotniy/elastic-image/internal/httpclient/mocks"
	"github.com/pzabolotniy/elastic-image/internal/logging"
	"github.com/pzabolotniy/monkey"
	"github.com/stretchr/testify/assert"
)

func TestNewFetchParams(t *testing.T) {
	testName := t.Name()
	testConfig := config.AppConfig{
		ServerConfig: &config.ServerConfig{
			Bind: "",
		},
		ImageConfig: &config.ImageConfig{
			CacheTTL:     3600 * time.Second,
			FetchTimeout: 10 * time.Second,
		},
	}
	timeout := testConfig.ImageConfig.FetchTimeout
	url := "http://any-host.local/image.jpeg"

	params := NewFetchParams(timeout, url)
	expected := &Params{
		Timeout: timeout,
		URL:     url,
	}
	assert.Equalf(t, expected, params, "%s - result ok", testName)
}

func TestGetImage(t *testing.T) {
	ctx := context.Background()
	testName := t.Name()
	testConfig := config.AppConfig{
		ServerConfig: &config.ServerConfig{
			Bind: "",
		},
		ImageConfig: &config.ImageConfig{
			CacheTTL:     3600 * time.Second,
			FetchTimeout: 10 * time.Second,
		},
	}
	testLogger := logging.GetLogger()
	ctx = logging.WithContext(ctx, testLogger)
	timeout := testConfig.ImageConfig.FetchTimeout
	url := "http://any-host.local/image.jpeg"
	var nilError error

	mockedBrowser := &mocks.Browser{}
	mockedResponse := &mocks.Responser{}
	monkey.Patch(httpclient.NewHTTPClient, func(logger logging.Logger, timeout time.Duration) httpclient.Browser {
		return mockedBrowser
	})

	mockedImage := []byte("i.am.image")
	mockedBrowser.On("Get", url).Return(mockedResponse, nilError)
	mockedResponse.On("RawBody").Return(mockedImage)
	params := NewFetchParams(timeout, url)
	response, err := GetImage(ctx, params)

	assert.NoErrorf(t, err, "%s - no errors, ok", testName)
	imageBytes, err := ioutil.ReadAll(response)
	if err != nil {
		panic(err)
	}
	assert.Equalf(t, mockedImage, imageBytes, "%s - image bytes ok", testName)

	mockedBrowser.AssertExpectations(t)
	mockedResponse.AssertExpectations(t)
}
