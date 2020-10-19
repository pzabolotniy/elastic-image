// Package fetch contains some funcs
// to download media data
package fetch

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/pzabolotniy/elastic-image/internal/httpclient"
	"github.com/pzabolotniy/elastic-image/internal/logging"
)

// Params is a container for image download parameters
type Params struct {
	Timeout time.Duration
	URL     string
}

// NewFetchParams is a constructor for Params
func NewFetchParams(timeout time.Duration, url string) *Params {
	params := &Params{
		Timeout: timeout,
		URL:     url,
	}
	return params
}

// GetImage downloads image by *Params
func GetImage(ctx context.Context, params *Params) (io.Reader, error) {
	logger := logging.FromContext(ctx)
	timeout := params.Timeout
	URL := params.URL

	httpClient := httpclient.NewHTTPClient(timeout)
	response, err := httpClient.Get(ctx, URL)
	if err != nil {
		logger.WithError(err).WithField("url", URL).Error("fetch image failed")
		return nil, err
	}

	imageReader := bytes.NewReader(response.RawBody)
	return imageReader, nil
}
