package fetch

import (
	"bytes"
	"github.com/pzabolotniy/elastic-image/internal/httpclient"
	"github.com/pzabolotniy/elastic-image/internal/logging"
	"io"
	"time"
)

type FetchParams struct {
	Timeout time.Duration
	URL string
	Logger logging.Logger
}

func NewFetchParams( timeout time.Duration, url string, logger logging.Logger ) *FetchParams {
	params := &FetchParams{
		Timeout:timeout,
		URL:url,
		Logger:logger,
	}
	return params
}

func GetImage( params *FetchParams ) ( io.Reader, error ) {
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
