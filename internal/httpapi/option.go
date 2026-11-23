package httpapi

import (
	"log/slog"
	"net/http"
)

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient sets the HTTP client used for API requests. The default
// client times out each request after 30 seconds.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient != nil {
			c.http = httpClient
		}
	}
}

// WithLogger sets the logger for request diagnostics. Headers, query
// parameters, and bodies are never logged because they carry the API key and
// solutions.
func WithLogger(logger *slog.Logger) Option {
	return func(c *Client) {
		if logger != nil {
			c.logger = logger
		}
	}
}

// WithBaseURL overrides the API root, for example to use a mirror or a
// local stub.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		if baseURL != "" {
			c.config.BaseURL = baseURL
		}
	}
}
