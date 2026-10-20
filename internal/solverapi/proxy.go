package solverapi

import (
	"fmt"

	"github.com/aarock1234/unicap/v3"
)

// ProxyFields holds the discrete proxy parameters most providers accept.
// Embed it in a provider task struct; the zero value is omitted from JSON.
type ProxyFields struct {
	ProxyType     string `json:"proxyType,omitempty"`
	ProxyAddress  string `json:"proxyAddress,omitempty"`
	ProxyPort     int    `json:"proxyPort,omitempty"`
	ProxyLogin    string `json:"proxyLogin,omitempty"`
	ProxyPassword string `json:"proxyPassword,omitempty"`
}

// NewProxyFields converts a unicap proxy. A nil proxy yields the zero value.
func NewProxyFields(proxy *unicap.Proxy) ProxyFields {
	if proxy == nil {
		return ProxyFields{}
	}

	return ProxyFields{
		ProxyType:     string(proxy.Type),
		ProxyAddress:  proxy.Host,
		ProxyPort:     proxy.Port,
		ProxyLogin:    proxy.Username,
		ProxyPassword: proxy.Password,
	}
}

// ProxyURL formats a proxy as "scheme://user:pass@host:port" for providers
// that take a single proxy string. Credentials are omitted when absent, and
// a nil proxy yields an empty string.
func ProxyURL(proxy *unicap.Proxy) string {
	if proxy == nil {
		return ""
	}

	return proxy.URL().String()
}

// RequireIPProxy reports an error when the proxy host is not a literal IP
// address, for providers that reject hostnames.
func RequireIPProxy(proxy *unicap.Proxy) error {
	if proxy == nil || proxy.HostIsIP() {
		return nil
	}

	return fmt.Errorf("Proxy.Host must be an IP address for this provider: %w", unicap.ErrUnsupportedOption)
}
