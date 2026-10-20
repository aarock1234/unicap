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
)

// maxResponseBytes bounds a response body. Solutions can embed page HTML or
// screenshots, so the limit is generous.
const maxResponseBytes = 16 << 20

// Call posts request to the named API method and decodes the reply into
// response. A non-zero errorId becomes a *unicap.Error classified by the
// provider's error rules, whatever the HTTP status. Server errors and rate
// limits without an error envelope are returned as plain transport errors so
// pollers treat them as transient.
func (c *Client) Call(ctx context.Context, method string, request Request, response Response) error {
	request.authenticate(c.config.APIKey)

	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("%s %s: encoding request: %w", c.config.Name, method, err)
	}

	statusCode, responseBody, err := c.post(ctx, method, body)
	if err != nil {
		return fmt.Errorf("%s %s: %w", c.config.Name, method, err)
	}

	decodeErr := json.Unmarshal(responseBody, response)

	if envelope := response.envelope(); decodeErr == nil && envelope.ErrorID != 0 {
		return c.config.Errors.classify(c.config.Name, envelope)
	}

	if statusCode != http.StatusOK {
		return c.statusError(method, statusCode)
	}

	if decodeErr != nil {
		return fmt.Errorf("%s %s: decoding response: %w", c.config.Name, method, decodeErr)
	}

	return nil
}

// post sends one request and returns the status code and the bounded body.
func (c *Client) post(ctx context.Context, method string, body []byte) (int, []byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.config.BaseURL+"/"+method, bytes.NewReader(body))
	if err != nil {
		return 0, nil, fmt.Errorf("creating request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	started := time.Now()

	response, err := c.http.Do(request)
	if err != nil {
		return 0, nil, fmt.Errorf("sending request: %w", err)
	}
	defer func() { _ = response.Body.Close() }()

	responseBody, err := io.ReadAll(io.LimitReader(response.Body, maxResponseBytes+1))
	if err != nil {
		return 0, nil, fmt.Errorf("reading response: %w", err)
	}

	if len(responseBody) > maxResponseBytes {
		return 0, nil, fmt.Errorf("response exceeds %d bytes", maxResponseBytes)
	}

	c.logger.DebugContext(ctx, "provider request finished",
		slog.String("provider", c.config.Name),
		slog.String("method", method),
		slog.Int("status_code", response.StatusCode),
		slog.Duration("duration", time.Since(started)),
	)

	return response.StatusCode, responseBody, nil
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
