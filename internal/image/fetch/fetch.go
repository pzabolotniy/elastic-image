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

// DownloadState contains download image
// and list of clients, who requested it
type DownloadState struct {
	Data []byte
	Subs []chan error
}

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
func GetImage(ctx context.Context, sharedDownload map[string]*DownloadState, params *Params) (io.Reader, error) {
	logger := logging.FromContext(ctx)
	timeout := params.Timeout
	URL := params.URL
	var imageReader io.Reader

	if dnState, ok := sharedDownload[URL]; ok {
		logger.WithField("url", URL).Trace("is fetching by another client")
		errCh := make(chan error, 1)
		dnState.Subs = append(dnState.Subs, errCh)
		if err := <-errCh; err != nil {
			logger.WithError(err).WithField("url", URL).Trace("fetch failed")
			delete(sharedDownload, URL)
			return nil, err
		}
		imageReader = bytes.NewReader(dnState.Data)
		logger.WithField("url", URL).Trace("fetched")
		delete(sharedDownload, URL)
	} else {
		subscribers := make([]chan error, 0, 1)
		downloadState := &DownloadState{
			Data: nil,
			Subs: subscribers,
		}
		sharedDownload[URL] = downloadState
		httpClient := httpclient.NewHTTPClient(timeout)
		response, err := httpClient.Get(ctx, URL)
		if err != nil {
			logger.WithError(err).WithField("url", URL).Error("fetch image failed")
			for _, subs := range downloadState.Subs {
				subs <- err
			}
			return nil, err
		}

		downloadState.Data = response.RawBody
		for _, subs := range downloadState.Subs {
			subs <- nil
		}
		imageReader = bytes.NewReader(response.RawBody)
	}

	return imageReader, nil
}
