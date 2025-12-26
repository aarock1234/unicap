package unicap

import (
	"log/slog"
)

// Option configures a Client
type Option func(*Client)

// WithLogger sets a custom logger
func WithLogger(logger *slog.Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
}

// WithPoller sets a custom poller
func WithPoller(poller *Poller) Option {
	return func(c *Client) {
		c.poller = poller
	}
}
