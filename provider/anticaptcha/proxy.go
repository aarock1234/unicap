package anticaptcha

import (
	"fmt"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/solverapi"
)

// proxyFields converts a proxy to Anti-Captcha's fields. Anti-Captcha takes
// IP addresses only and names HTTP and HTTPS proxies alike "http".
func proxyFields(proxy *unicap.Proxy) (solverapi.ProxyFields, error) {
	if err := solverapi.RequireIPProxy(proxy); err != nil {
		return solverapi.ProxyFields{}, err
	}

	fields := solverapi.NewProxyFields(proxy)
	if proxy != nil && proxy.Type == unicap.ProxyTypeHTTPS {
		fields.ProxyType = string(unicap.ProxyTypeHTTP)
	}

	return fields, nil
}

// requireUserAgentWithProxy reports an error when a proxy is set without the
// User-Agent Anti-Captcha needs to emulate the browser.
func requireUserAgentWithProxy(proxy *unicap.Proxy, userAgent string) error {
	if proxy == nil || userAgent != "" {
		return nil
	}

	return fmt.Errorf("UserAgent is required with a proxy: %w", unicap.ErrInvalidTask)
}
