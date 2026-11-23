package anticaptcha

import (
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

// userAgentWithProxyOnly enforces Anti-Captcha's User-Agent rule: proxy task
// types require one, and proxyless task types do not accept one.
func userAgentWithProxyOnly(proxy *unicap.Proxy, userAgent string) error {
	return solverapi.Check(
		solverapi.Unsupported(
			solverapi.Setting("UserAgent without Proxy", proxy == nil && userAgent != ""),
		),
		solverapi.Required(
			solverapi.Setting("UserAgent with a Proxy", proxy == nil || userAgent != ""),
		),
	)
}
