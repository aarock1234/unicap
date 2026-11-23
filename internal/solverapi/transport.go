package solverapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/aarock1234/unicap/v3"
)

// maxResponseBytes bounds a response body. Solutions can embed page HTML or
// screenshots, so the limit is generous.
const maxResponseBytes = 16 << 20

// Call posts request to the named API method and decodes the reply into
// response. A non-zero errorId becomes a *unicap.Error classified by the
// provider's error rules, whatever the HTTP status. Server errors and rate
// limits without an error envelope are returned as plain transport errors so
// pollers treat them as transient. A request that cannot be encoded or built,
// and a response over the size limit, are final provider errors.
func (c *Client) Call(ctx context.Context, method string, request Request, response Response) error {
	request.authenticate(c.config.APIKey)

	body, err := json.Marshal(request)
	if err != nil {
		return c.providerError(
			"INVALID_REQUEST",
			fmt.Sprintf("%s: encoding request: %v", method, err),
			unicap.ErrInvalidTask,
		)
	}

	statusCode, responseBody, err := c.post(ctx, method, body)
	if err != nil {
		return err
	}

	var envelope Envelope
	if err := json.Unmarshal(responseBody, &envelope); err == nil && envelope.ErrorID != 0 {
		return c.config.Errors.classify(c.config.Name, &envelope)
	}

	if statusCode != http.StatusOK {
		return c.statusError(method, statusCode)
	}

	if err := json.Unmarshal(responseBody, response); err != nil {
		return fmt.Errorf("%s %s: decoding response: %w", c.config.Name, method, err)
	}

	return nil
}

// post sends one request and returns the status code and the bounded body.
// Failures to send the request or read the response are plain errors, which
// callers treat as transient; every other failure is a final provider error.
func (c *Client) post(ctx context.Context, method string, body []byte) (int, []byte, error) {
	request, err := c.newRequest(ctx, method, body)
	if err != nil {
		return 0, nil, err
	}

	started := time.Now()

	response, err := c.http.Do(request)
	if err != nil {
		return 0, nil, fmt.Errorf("%s %s: sending request: %w", c.config.Name, method, err)
	}
	defer func() { _ = response.Body.Close() }()

	responseBody, err := io.ReadAll(io.LimitReader(response.Body, maxResponseBytes+1))
	if err != nil {
		return 0, nil, fmt.Errorf("%s %s: reading response: %w", c.config.Name, method, err)
	}

	if len(responseBody) > maxResponseBytes {
		return 0, nil, c.providerError(
			"RESPONSE_TOO_LARGE",
			fmt.Sprintf("%s response exceeds %d bytes", method, maxResponseBytes),
			nil,
		)
	}

	c.logger.DebugContext(ctx, "provider request finished",
		slog.String("provider", c.config.Name),
		slog.String("method", method),
		slog.Int("status_code", response.StatusCode),
		slog.Duration("duration", time.Since(started)),
	)

	return response.StatusCode, responseBody, nil
}

// newRequest builds the JSON POST request for the named API method. A request
// that cannot be built is a final provider error.
func (c *Client) newRequest(ctx context.Context, method string, body []byte) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.config.BaseURL+"/"+method, bytes.NewReader(body))
	if err != nil {
		return nil, c.providerError(
			"INVALID_REQUEST",
			fmt.Sprintf("%s: creating request: %v", method, err),
			nil,
		)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	return request, nil
}

// statusError reports an HTTP failure without an error envelope. Server
// errors and rate limits stay plain errors, which callers treat as
// transient; other statuses become final provider errors.
func (c *Client) statusError(method string, statusCode int) error {
	if statusCode >= http.StatusInternalServerError || statusCode == http.StatusTooManyRequests {
		return fmt.Errorf("%s %s: unexpected http status %d", c.config.Name, method, statusCode)
	}

	return c.providerError(
		"HTTP_"+strconv.Itoa(statusCode),
		fmt.Sprintf("%s returned http status %d", method, statusCode),
		nil,
	)
}
