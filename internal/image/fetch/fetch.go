package fetch

import (
	"bytes"
	"github.com/pzabolotniy/elastic-image/internal/httpclient"
	"github.com/pzabolotniy/elastic-image/internal/logging"
	"io"
	"time"
)

// Params is a container for image download parameters
type Params struct {
	Timeout time.Duration
	URL     string
	Logger  logging.Logger
}

// NewFetchParams is a constructor for Params
func NewFetchParams(timeout time.Duration, url string, logger logging.Logger) *Params {
	params := &Params{
		Timeout: timeout,
		URL:     url,
		Logger:  logger,
	}
	return params
}

// GetImage downloads image by *Params
func GetImage(params *Params) (io.Reader, error) {
	logger := params.Logger
	timeout := params.Timeout
	URL := params.URL

	httpClient := httpclient.NewHTTPClient(logger, timeout)
	response, err := httpClient.Get(URL)
	if err != nil {
		logger.Errorf("get url '%s' failed: '%s'", URL, err)
		return nil, err
	}

	imageBin := response.RawBody()
	imageReader := bytes.NewReader(imageBin)
	return imageReader, nil
}
