package unicap

import "log/slog"

// Option configures a Client.
type Option func(*Client)

// WithLogger sets a custom logger. The logger is also used by the default
// poller unless a custom poller is supplied with WithPoller.
func WithLogger(logger *slog.Logger) Option {
	return func(c *Client) {
		if logger != nil {
			c.logger = logger
		}
	}
}

// WithPoller sets a custom poller, overriding the default.
func WithPoller(poller *Poller) Option {
	return func(c *Client) {
		if poller != nil {
			c.poller = poller
		}
	}
}
