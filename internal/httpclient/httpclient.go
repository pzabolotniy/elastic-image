package httpclient

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/pzabolotniy/elastic-image/internal/logging"
	"gopkg.in/resty.v1"
)

// Client is a container for http client parameters
type Client struct {
	client  *resty.Client
	timeout time.Duration
}

// ResponseContainer contains http response data
type ResponseContainer struct {
	RawBody    []byte
	StatusCode int
	StatusLine string
}

// NewHTTPClient is a constructor for Browser
func NewHTTPClient(timeout time.Duration) *Client {
	restClient := resty.New()

	env := &Client{
		client:  restClient,
		timeout: timeout,
	}

	return env
}

// Get implements http GET method
func (c *Client) Get(ctx context.Context, url string) (*ResponseContainer, error) {
	return c.ExecuteRequest(ctx, http.MethodGet, url)
}

// ExecuteRequest implements http request with any method
func (c *Client) ExecuteRequest(ctx context.Context, method, url string) (*ResponseContainer, error) {
	logger := logging.FromContext(ctx)
	timeout := c.timeout

	if url == "" {
		logger.Warn("Empty URL of the remote service")
		return nil, errors.New("got empty URL of the remote service")
	}

	var httpResponse *resty.Response
	var err error
	rClient := c.client

	rClient.SetTimeout(timeout)
	request := rClient.R()
	switch method {
	case http.MethodGet:
		logger.WithFields(logging.Fields{
			"http_method": method,
			"url":         url,
		}).Trace("http request")
		httpResponse, err = request.Get(url)
		if err != nil {
			logger.WithError(err).Error("request failed")
			return nil, errors.New("request failed")
		}
	default:
		logger.WithField("http_method", method).Warn("http method is not implemented")
		return nil, errors.New("method is not implemented")
	}

	statusLine := httpResponse.Status()
	statuscode := httpResponse.StatusCode()
	responseSize := httpResponse.Size()
	logger.WithField("status_line", statusLine).Trace("response status")
	logger.WithField("response_size", responseSize).Trace("response size in bytes")
	var responseBody []byte
	if responseSize > 0 {
		responseBody = httpResponse.Body()
	}

	restResponse := NewResponse(statuscode, statusLine, responseBody)

	if !restResponse.IsSuccess() {
		err = errors.New(statusLine)
	}

	return restResponse, err
}

// NewResponse is constructor for Responser interface
// contains rest-call response code and body
func NewResponse(statusCode int, statusLine string, body []byte) *ResponseContainer {
	restResponse := &ResponseContainer{
		StatusCode: statusCode,
		RawBody:    body,
		StatusLine: statusLine,
	}
	return restResponse
}

// IsSuccess returns true if http code is any if 2XX
// otherwise returns false
func (resp *ResponseContainer) IsSuccess() bool {
	success := false
	switch resp.StatusCode {
	case http.StatusOK:
		success = true
	case http.StatusCreated:
		success = true
	case http.StatusAccepted:
		success = true
	case http.StatusNoContent:
		success = true
	}
	return success
}
