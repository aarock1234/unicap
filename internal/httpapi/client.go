// Package httpapi sends bounded HTTP requests for provider clients and
// classifies failures as transient transport errors or final provider errors.
// It knows nothing about any provider's protocol.
package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/aarock1234/unicap/v3"
)

// defaultHTTPTimeout bounds a single API request.
const defaultHTTPTimeout = 30 * time.Second

// maxResponseBytes bounds a response body. Solutions can embed page HTML or
// screenshots, so the limit is generous.
const maxResponseBytes = 16 << 20

// jsonContentType is the media type of JSON request and response bodies.
const jsonContentType = "application/json"

// Config describes the API a Client talks to.
type Config struct {
	// Name is the provider identifier reported in errors and logs.
	Name string

	// BaseURL is the API root, without a trailing slash.
	BaseURL string

	// Header is sent with every request, for example an Authorization
	// header. It is never logged.
	Header http.Header
}

// Request is one HTTP exchange.
type Request struct {
	// Method is the HTTP method.
	Method string

	// Path is appended to the base URL and starts with a slash.
	Path string

	// Query holds optional URL query parameters.
	Query url.Values

	// Header holds optional per-request headers.
	Header http.Header

	// Body is JSON-encoded and sent as application/json when not nil.
	Body any

	// Encoded is a pre-encoded body, such as a multipart form. It is used
	// instead of Body when set.
	Encoded *EncodedBody
}

// operation names the exchange in errors, as the HTTP method and the path.
func (r Request) operation() string {
	return r.Method + " " + r.Path
}

// EncodedBody is a request body the caller encoded itself.
type EncodedBody struct {
	// ContentType is the Content-Type header value.
	ContentType string

	// Data is the body.
	Data []byte
}

// Response is a received HTTP response with its bounded body.
type Response struct {
	// StatusCode is the HTTP status code.
	StatusCode int

	// Header holds the response headers.
	Header http.Header

	// Body is the response body.
	Body []byte
}

// Client sends HTTP requests to one provider API. It is safe for concurrent
// use.
type Client struct {
	config Config
	http   *http.Client
	logger *slog.Logger
}

// New creates a client for the API described by config.
func New(config Config, opts ...Option) *Client {
	c := &Client{
		config: config,
		http: &http.Client{
			Timeout: defaultHTTPTimeout,
		},
		logger: slog.New(slog.DiscardHandler),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Name returns the provider identifier.
func (c *Client) Name() string {
	return c.config.Name
}

// Do sends request and returns the response with its bounded body, whatever
// the HTTP status. Failures to send the request or read the response are
// plain errors, which callers treat as transient. A request that cannot be
// encoded or built, and a response over the size limit, are final provider
// errors.
func (c *Client) Do(ctx context.Context, request Request) (*Response, error) {
	body, err := c.encodeBody(request)
	if err != nil {
		return nil, err
	}

	httpRequest, err := c.newRequest(ctx, request, body)
	if err != nil {
		return nil, err
	}

	started := time.Now()

	httpResponse, err := c.http.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("%s %s: sending request: %w", c.config.Name, request.operation(), err)
	}
	defer func() { _ = httpResponse.Body.Close() }()

	responseBody, err := c.readBody(request, httpResponse.Body)
	if err != nil {
		return nil, err
	}

	c.logger.DebugContext(ctx, "provider request finished",
		slog.String("provider", c.config.Name),
		slog.String("method", request.Method),
		slog.String("path", request.Path),
		slog.Int("status_code", httpResponse.StatusCode),
		slog.Duration("duration", time.Since(started)),
	)

	return &Response{
		StatusCode: httpResponse.StatusCode,
		Header:     httpResponse.Header,
		Body:       responseBody,
	}, nil
}

// Error builds a final provider error that the provider did not report as an
// error code.
func (c *Client) Error(code, message string, sentinel error) *unicap.Error {
	return &unicap.Error{
		Provider: c.config.Name,
		Code:     code,
		Message:  message,
		Err:      sentinel,
	}
}

// StatusError reports an HTTP failure the provider did not explain. Server
// errors and rate limits stay plain errors, which callers treat as
// transient; other statuses become final provider errors.
func (c *Client) StatusError(operation string, statusCode int) error {
	if statusCode >= http.StatusInternalServerError || statusCode == http.StatusTooManyRequests {
		return fmt.Errorf("%s %s: unexpected http status %d", c.config.Name, operation, statusCode)
	}

	return c.Error(
		"HTTP_"+strconv.Itoa(statusCode),
		fmt.Sprintf("%s returned http status %d", operation, statusCode),
		nil,
	)
}

// encodeBody returns the body to send, or nil when the request has none. A
// pre-encoded body is used as is; any other body is encoded as JSON, and one
// that cannot be encoded is a final provider error.
func (c *Client) encodeBody(request Request) (*EncodedBody, error) {
	if request.Encoded != nil {
		return request.Encoded, nil
	}

	if request.Body == nil {
		return nil, nil
	}

	data, err := json.Marshal(request.Body)
	if err != nil {
		return nil, c.Error(
			"INVALID_REQUEST",
			fmt.Sprintf("%s: encoding request: %v", request.operation(), err),
			unicap.ErrInvalidTask,
		)
	}

	return &EncodedBody{
		ContentType: jsonContentType,
		Data:        data,
	}, nil
}

// newRequest builds the HTTP request with the client's headers, then the
// request's own. A request that cannot be built is a final provider error.
func (c *Client) newRequest(ctx context.Context, request Request, body *EncodedBody) (*http.Request, error) {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body.Data)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, request.Method, c.requestURL(request), reader)
	if err != nil {
		return nil, c.Error(
			"INVALID_REQUEST",
			fmt.Sprintf("%s: creating request: %v", request.operation(), err),
			nil,
		)
	}

	setHeader(httpRequest.Header, c.config.Header)
	setHeader(httpRequest.Header, request.Header)
	httpRequest.Header.Set("Accept", jsonContentType)

	if body != nil {
		httpRequest.Header.Set("Content-Type", body.ContentType)
	}

	return httpRequest, nil
}

// requestURL joins the base URL, the request path, and the encoded query.
func (c *Client) requestURL(request Request) string {
	requestURL := c.config.BaseURL + request.Path

	if len(request.Query) == 0 {
		return requestURL
	}

	return requestURL + "?" + request.Query.Encode()
}

// readBody reads the response body up to the size limit. A read failure is a
// plain error; a body over the limit is a final provider error.
func (c *Client) readBody(request Request, body io.Reader) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(body, maxResponseBytes+1))
	if err != nil {
		return nil, fmt.Errorf("%s %s: reading response: %w", c.config.Name, request.operation(), err)
	}

	if len(data) > maxResponseBytes {
		return nil, c.Error(
			"RESPONSE_TOO_LARGE",
			fmt.Sprintf("%s response exceeds %d bytes", request.operation(), maxResponseBytes),
			nil,
		)
	}

	return data, nil
}

// setHeader replaces the values in target with those in source, key by key.
func setHeader(target, source http.Header) {
	for key, values := range source {
		target.Del(key)

		for _, value := range values {
			target.Add(key, value)
		}
	}
}
