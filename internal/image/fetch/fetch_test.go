package fetch

import (
	"bou.ke/monkey"
	"github.com/pzabolotniy/elastic-image/internal/config"
	"github.com/pzabolotniy/elastic-image/internal/httpclient"
	"github.com/pzabolotniy/elastic-image/internal/httpclient/mocks"
	"github.com/pzabolotniy/elastic-image/internal/logging"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"time"
)

func TestNewFetchParams(t *testing.T) {
	testName := t.Name()
	testConfig := config.GetConfig()
	testLogger := testConfig.APILogger
	timeout := testConfig.Timeout
	url := "http://any-host.local/image.jpeg"

	params := NewFetchParams(timeout, url, testLogger)
	expected := &Params{
		Timeout: timeout,
		Logger:  testLogger,
		URL:     url,
	}
	assert.Equalf(t, expected, params, "%s - result ok", testName)
}

func TestGetImage(t *testing.T) {
	testName := t.Name()
	testConfig := config.GetConfig()
	testLogger := testConfig.APILogger
	timeout := testConfig.Timeout
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
	params := NewFetchParams(timeout, url, testLogger)
	response, err := GetImage(params)

	assert.NoErrorf(t, err, "%s - no errors, ok", testName)
	imageBytes, err := ioutil.ReadAll(response)
	if err != nil {
		panic(err)
	}
	assert.Equalf(t, mockedImage, imageBytes, "%s - image bytes ok", testName)

	mockedBrowser.AssertExpectations(t)
	mockedResponse.AssertExpectations(t)
}
