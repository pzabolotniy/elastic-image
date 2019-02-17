package httpclient

import (
	"errors"
	"github.com/pzabolotniy/elastic-image/internal/logging"
	"gopkg.in/resty.v1"
	"net/http"
	"time"
)

// list of implemented HTTP methods
const (
	methodGet = "GET"
)

// Env is a container for http client parameters
// must implement Browser interface
type Env struct {
	client  *resty.Client
	timeout time.Duration
	logger  logging.Logger
}

// Browser describes methods
// which can be made by http client
type Browser interface {
	Logger() logging.Logger
	ExecuteRequest(method string, url string) (Responser, error)
	Get(url string) (Responser, error)
}

// ResponseContainer contains http response data
// must implement Responser
type ResponseContainer struct {
	rawBody    []byte
	statusCode int
	statusLine string
}

// Responser contains getters to http response data
type Responser interface {
	RawBody() []byte
	StatusCode() int
	IsSuccess() bool
}

// NewHTTPClient is a constructor for Browser
func NewHTTPClient(logger logging.Logger, timeout time.Duration) Browser {
	restClient := resty.New()

	env := &Env{
		client:  restClient,
		timeout: timeout,
		logger:  logger,
	}
	var _ Browser = env

	return env
}

// Timeout is a getter for Env.timeout
func (client *Env) Timeout() time.Duration {
	return client.timeout
}

// Logger is a getter for Env.logger
func (client *Env) Logger() logging.Logger {
	return client.logger
}

// Get implements http GET method
func (client *Env) Get(url string) (Responser, error) {
	return client.ExecuteRequest(methodGet, url)
}

// ExecuteRequest implements http request with any method
func (client *Env) ExecuteRequest(method string, url string) (Responser, error) {
	logger := client.Logger()
	timeout := client.Timeout()

	if len(url) == 0 {
		logger.Warnf("Empty URL of the remote service")
		return nil, errors.New("Got empty URL of the remote service")
	}

	var httpServiceResponse *resty.Response
	var err error
	rClient := client.client

	rClient.SetTimeout(time.Duration(timeout) * time.Millisecond)
	request := rClient.R()
	switch method {
	case methodGet:
		logger.Tracef("request: %s %s", methodGet, url)
		httpServiceResponse, err = request.Get(url)
	default:
		logger.Warnf("HTTP request method %s is not implemented", method)
		return nil, errors.New("method is not implemented")
	}

	statusLine := httpServiceResponse.Status()
	statuscode := httpServiceResponse.StatusCode()
	responseSize := httpServiceResponse.Size()
	logger.Tracef("response status: %s", statusLine)
	logger.Tracef("response size in bytes: %d", responseSize)
	var responseBytes []byte
	if responseSize > 0 {
		responseBytes = httpServiceResponse.Body()
	}

	restResponse := NewResponse(statuscode, statusLine, responseBytes)

	if err == nil {
		if !restResponse.IsSuccess() {
			err = errors.New(statusLine)
		}
	} else {
		logger.Warnf("destination is unavailable: '%s'", err.Error())
	}

	return restResponse, err
}

// NewResponse is constructor for Responser interface
// contains rest-call response code and body
func NewResponse(statusCode int, statusLine string, body []byte) Responser {
	restResponse := &ResponseContainer{
		statusCode: statusCode,
		rawBody:    body,
		statusLine: statusLine,
	}
	var _ Responser = restResponse
	return restResponse
}

// RawBody is a getter for ResponseContainer.rawBody
func (resp *ResponseContainer) RawBody() []byte {
	return resp.rawBody
}

// StatusCode is a getter for ResponseContainer.statusCode
func (resp *ResponseContainer) StatusCode() int {
	return resp.statusCode
}

// IsSuccess returns true if http code is any if 2XX
// otherwise returns false
func (resp *ResponseContainer) IsSuccess() bool {
	success := false
	switch resp.statusCode {
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