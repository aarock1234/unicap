package captchaai

import (
	"net"
	"strconv"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/solverapi"
)

// proxyTypes are the proxy types CaptchaAI accepts, each with the value its
// proxytype field takes. It accepts no SOCKS proxy.
var proxyTypes = map[unicap.ProxyType]string{
	unicap.ProxyTypeHTTP:  "HTTP",
	unicap.ProxyTypeHTTPS: "HTTPS",
}

// proxyFields returns the two form fields that describe an optional proxy:
// proxy, which reads "login:password@host:port", and proxytype. A nil proxy
// yields no field. CaptchaAI takes HTTP and HTTPS proxies only and requires
// credentials, so a proxy of another type and one without a username are
// reported as unsupported.
func proxyFields(source *unicap.Proxy) (form, error) {
	if source == nil {
		return nil, nil
	}

	proxyType, accepted := proxyTypes[source.Type]

	if err := solverapi.Unsupported(
		solverapi.Setting("Proxy with type "+string(source.Type), !accepted),
		solverapi.Setting("Proxy without a username", !source.HasCredentials()),
	); err != nil {
		return nil, err
	}

	address := net.JoinHostPort(source.Host, strconv.Itoa(source.Port))

	return form{}.
		with("proxy", source.Username+":"+source.Password+"@"+address).
		with("proxytype", proxyType), nil
}
